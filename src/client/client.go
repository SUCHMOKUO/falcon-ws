package client

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"socks5"
	"sync"
	"time"
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

	// websocket dialer.
	dialer = websocket.Dialer{
		HandshakeTimeout: timeout,
		ReadBufferSize:   bufSize,
		WriteBufferSize:  bufSize,
		WriteBufferPool:  new(sync.Pool),
	}
)

// NewClient create a client.
func NewClient(socks5Addr string, serverAddr string) {
	socks5.ListenAndServe(socks5Addr, func(c net.Conn, t *socks5.Target) {
		// url encode.
		host := base64.URLEncoding.EncodeToString([]byte(t.Host))
		port := base64.URLEncoding.EncodeToString([]byte(t.Port))
		url := fmt.Sprintf("ws://%s/free?h=%s&p=%s", serverAddr, host, port)
		// log.Println(url)
		ws, res, err := dialer.Dial(url, nil)
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
