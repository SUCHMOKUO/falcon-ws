package client

import (
	"crypto/tls"
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/SUCHMOKUO/falcon-ws/socks5"
	"github.com/SUCHMOKUO/falcon-ws/util"
	"github.com/gorilla/websocket"
)

const (
	// transmission type
	dataT = websocket.BinaryMessage

	// connection time out
	timeout = 10 * time.Second

	// buffer size.
	bufSize = 1024
)

var (
	// buffer pool.
	bufPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, bufSize)
		},
	}

	// fake User-Agent.
	defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML"
)

// Client
type Client struct {
	// local socks5 address for listening.
	Socks5Addr string

	// server info.
	Host string
	Port string
	wsAddr string

	// fake 'Host' field in request header
	// for against qos.
	FakeHost string

	// fake User-Agent.
	UserAgent string

	// secure flag. true for using https
	// instead of http.
	Secure bool

	// lookup flag. if it's false, client will
	// not lookup the server ip and cache it.
	Lookup bool

	// ipv6 flag. if it's true, the ipv6 address
	// of proxy server will be used first.
	IPv6 bool

	// http header
	header http.Header

	// websocket dialer.
	dialer websocket.Dialer
}

// NewClient return new falcon client instance.
func NewClient(socks5addr, serveraddr, fakehost, useragent string, secure, lookup, ipv6 bool) *Client {
	c := new(Client)
	c.Socks5Addr = socks5addr
	c.FakeHost = fakehost
	c.UserAgent = useragent
	c.Secure = secure
	c.Lookup = lookup
	c.IPv6 = ipv6
	c.dialer = websocket.Dialer{
		HandshakeTimeout: timeout,
		ReadBufferSize:   bufSize,
		WriteBufferSize:  bufSize,
		WriteBufferPool:  new(sync.Pool),
	}

	host, port, err := net.SplitHostPort(serveraddr)
	if err != nil {
		log.Fatalln("server address format error!")
	}
	c.Host = host
	c.Port = port

	// tls setting.
	if c.Secure {
		tlsCfg := new(tls.Config)
		tlsCfg.ServerName = c.Host
		c.dialer.TLSClientConfig = tlsCfg
	}

	// init fake header.
	reqHeader := http.Header{}
	if c.FakeHost != "" && !c.Secure {
		// add fake host field.
		reqHeader.Set("Host", fakehost)
	} else {
		reqHeader.Set("Host", c.Host)
	}
	// fake user-agent field.
	if c.UserAgent != "" {
		reqHeader.Set("User-Agent", c.UserAgent)
	} else {
		reqHeader.Set("User-Agent", defaultUserAgent)
	}
	c.header = reqHeader

	// lookup.
	var ip net.IP
	if c.Lookup && util.IsDomain(c.Host) {
		log.Println("lookup for server", c.Host, "ip address...")
		ip, err = util.Lookup(c.Host, c.IPv6)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// get ws address.
	var schema string
	if secure {
		schema = "wss://"
	}
	if !secure {
		schema = "ws://"
	}
	if !lookup {
		c.wsAddr = schema + c.Host + ":" + c.Port
		return c
	}
	var ipstr string
	if util.IsIPv4(ip) {
		ipstr = ip.String()
	}
	if util.IsIPv6(ip) {
		ipstr = "[" + ip.String() + "]"
	}
	c.wsAddr = schema + ipstr + ":" + c.Port
	return c
}

// Run falcon.
func (c *Client) Run() {
	log.Println("falcon server:", c.wsAddr)
	log.Println("use host:", c.header.Get("Host"))
	log.Println("use user-agent:", c.header.Get("User-Agent"))

	// start socks5 server.
	socks5.ListenAndServe(c.Socks5Addr, func(conn net.Conn, t *socks5.Target) {
		// url encode.
		host := base64.URLEncoding.EncodeToString([]byte(t.Host))
		port := base64.URLEncoding.EncodeToString([]byte(t.Port))
		url := c.wsAddr + "/free?h=" + host + "&p=" + port
		ws, res, err := c.dialer.Dial(url, c.header)
		if err != nil {
			log.Println("Dial proxy server error:", err)
			conn.Close()
			return
		}
		if res.StatusCode == 404 {
			log.Println("Empty target address or port.")
			conn.Close()
			return
		}

		go send(conn, ws)
		go recive(ws, conn)
	})
}

func send(conn net.Conn, ws *websocket.Conn) {
	defer ws.Close()
	defer conn.Close()

	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		if n > 0 {
			err = ws.WriteMessage(dataT, buf[:n])
			if err != nil {
				return
			}
		}
	}
}

func recive(ws *websocket.Conn, conn net.Conn) {
	defer ws.Close()
	defer conn.Close()

	for {
		msgT, data, err := ws.ReadMessage()

		if msgT != dataT || err != nil {
			return
		}

		_, err = conn.Write(data)

		if err != nil {
			return
		}
	}
}
