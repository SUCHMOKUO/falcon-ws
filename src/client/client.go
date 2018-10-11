package client

import (
	"encoding/base64"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"socks5"
	"strings"
	"sync"
	"time"
	"util"
)

const (
	// transmission type
	dataT = websocket.BinaryMessage

	// connection time out
	timeout = 10 * time.Second

	// buffer size.
	bufSize = 1024

	// fake User-Agent.
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML"
)

var (
	// buffer pool.
	bufPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, bufSize)
		},
	}

	// websocket dialer.
	dialer = websocket.Dialer{
		HandshakeTimeout: timeout,
		ReadBufferSize:   bufSize,
		WriteBufferSize:  bufSize,
		WriteBufferPool:  new(sync.Pool),
	}
)

// Config for client.
type Config struct {
	// local socks5 address for listening.
	Socks5Addr string

	// proxy server address.
	ServerAddr string

	// fake 'Host' field in request header
	// for against qos.
	FakeHost string

	// secure flag. true for using https
	// instead of http.
	Secure bool

	// lookup flag. if it's false, client will
	// not lookup the server ip and cache it.
	Lookup bool
}

func (conf *Config) lookupServer() {
	if conf.Secure || !conf.Lookup {
		// https enabled or non-lookup
		// flag is true.
		return
	}

	addr := strings.Split(conf.ServerAddr, ":")
	if len(addr) < 2 {
		log.Fatalln("server address format error!")
	}

	host := addr[0]
	port := addr[1]

	ok, err := util.IsDomain(host)
	if err != nil {
		log.Fatalln("server address format error!")
	}

	if !ok {
		// is ip.
		return
	}

	// is domain.
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Fatalln(err)
	}

	ip := ips[0]

	if ip.To4() != nil {
		// ipv4.
		conf.ServerAddr = ip.String() + ":" + port
	} else {
		// ipv6.
		conf.ServerAddr = "[" + ip.String() + "]:" + port
	}
}

// NewClient create a client.
func NewClient(conf *Config) {
	log.Println("Client initializing, please wait...")
	conf.lookupServer()
	log.Println("Done! get server address:", conf.ServerAddr)

	// request header.
	reqHeader := http.Header{}
	if conf.FakeHost != "" {
		// add fake host field.
		reqHeader.Set("Host", conf.FakeHost)
	}
	// fake user-agent field.
	reqHeader.Set("User-Agent", userAgent)

	// start socks5 server.
	log.Println("Socks5 server listening at", conf.Socks5Addr)
	socks5.ListenAndServe(conf.Socks5Addr, func(c net.Conn, t *socks5.Target) {
		// url encode.
		host := base64.URLEncoding.EncodeToString([]byte(t.Host))
		port := base64.URLEncoding.EncodeToString([]byte(t.Port))
		url := conf.ServerAddr + "/free?h=" + host + "&p=" + port;

		if conf.Secure {
			// https.
			url = "wss://" + url
		} else {
			// http.
			url = "ws://" + url
		}

		ws, res, err := dialer.Dial(url, reqHeader)
		if err != nil {
			log.Println("Dial proxy server error:", err)
			c.Close()
			return
		}
		if res.StatusCode == 404 {
			log.Println("Empty target address or port.")
			c.Close()
			return
		}

		go send(c, ws)
		go recive(ws, c)
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
