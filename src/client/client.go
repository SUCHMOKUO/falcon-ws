package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"socks5"
	"sync"
	"time"
	"net/http"
	"encoding/base64"
)

const (
	// transmission type
	dataT = websocket.BinaryMessage

	// connection time out
	timeout = 10 * time.Second

	// buffer size.
	bufSize = 1024

	// fake User-Agent.
	userAgent = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.2 (KHTML, like Gecko) Chrome/22.0.1216.0 Safari/537.2"
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
}

// NewClient create a client.
func NewClient(conf *Config) {
	// request header.
	reqHeader := http.Header{}
	if conf.FakeHost != "" {
		// add fake host field.
		reqHeader.Set("Host", conf.FakeHost)
	}
	// fake user-agent field.
	reqHeader.Set("User-Agent", userAgent)

	// start socks5 server.
	socks5.ListenAndServe(conf.Socks5Addr, func(c net.Conn, t *socks5.Target) {
		// url encode.
		host := base64.URLEncoding.EncodeToString([]byte(t.Host))
		port := base64.URLEncoding.EncodeToString([]byte(t.Port))
		url := fmt.Sprintf("%s/free?h=%s&p=%s", conf.ServerAddr, host, port)
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
