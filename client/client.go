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
	// connection time out
	timeout = 10 * time.Second

	// buffer size.
	bufSize = 1024
)

var (
	// fake User-Agent.
	defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML"
)

// Client
type Client struct {
	// local socks5 address for listening.
	Socks5Addr string

	// server info.
	Host   string
	Port   string
	WSAddr string

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
	Header http.Header

	// websocket dialer.
	Dialer websocket.Dialer
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
	c.Dialer = websocket.Dialer{
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
		c.Dialer.TLSClientConfig = tlsCfg
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
	c.Header = reqHeader

	// get ws address.
	var schema string
	if secure {
		schema = "wss://"
	} else {
		schema = "ws://"
	}
	if !lookup || !util.IsDomain(c.Host) {
		c.WSAddr = schema + c.Host + ":" + c.Port
		return c
	}
	// lookup.
	log.Println("lookup for server", c.Host, "ip address...")
	ip, err := util.Lookup(c.Host, c.IPv6)
	if err != nil {
		log.Fatalln(err)
	}
	var ipstr string
	if util.IsIPv4(ip) {
		ipstr = ip.String()
	}
	if util.IsIPv6(ip) {
		ipstr = "[" + ip.String() + "]"
	}
	c.WSAddr = schema + ipstr + ":" + c.Port
	return c
}

// Run falcon.
func (c *Client) Run() {
	log.Println("falcon server:", c.WSAddr)
	log.Println("use host:", c.Header.Get("Host"))
	log.Println("use user-agent:", c.Header.Get("User-Agent"))

	// start socks5 server.
	socks5.ListenAndServe(c.Socks5Addr, func(conn net.Conn, t *socks5.Target) {
		// url encode.
		host := base64.URLEncoding.EncodeToString([]byte(t.Host))
		port := base64.URLEncoding.EncodeToString([]byte(t.Port))
		url := c.WSAddr + "/free?h=" + host + "&p=" + port
		ws, res, err := c.Dialer.Dial(url, c.Header)
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

		go util.WSToConn(conn, ws)
		go util.ConnToWS(ws, conn)
	})
}