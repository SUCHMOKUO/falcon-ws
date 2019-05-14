package server

import (
	"errors"
	"github.com/SUCHMOKUO/falcon-ws/tunnel"
	"github.com/SUCHMOKUO/falcon-ws/util"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	// connection time out
	timeout = 10 * time.Second

	// buffer size.
	bufSize = 1500
)

var (
	// websocket upgrader.
	upgrader = websocket.Upgrader{
		HandshakeTimeout: timeout,
		ReadBufferSize:   bufSize,
		WriteBufferSize:  bufSize,
		WriteBufferPool:  new(sync.Pool),
	}

	errQueryParams = errors.New("query params error")
)

func handleProxyReq(w http.ResponseWriter, r *http.Request) {
	addr, err := getTarget(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	ch := make(chan net.Conn, 1)
	go connectTarget(addr, ch)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	t := &tunnel.Tunnel{*ws}

	conn := <-ch
	if conn == nil {
		t.Close()
		return
	}

	go util.Copy(conn, t)
	go util.Copy(t, conn)
}

func getTarget(r *http.Request) (string, error) {
	target := r.URL.Query()
	hostEncoded := target.Get("h")
	portEncoded := target.Get("p")

	var hasErr = false
	var host, port string
	var err error

	if hostEncoded == "" || portEncoded == "" {
		hasErr = true
	} else {
		// decode url.
		host, err = util.Decode(hostEncoded)
		if err != nil {
			hasErr = true
		}
		port, err = util.Decode(portEncoded)
		if err != nil {
			hasErr = true
		}
	}

	if hasErr {
		return "", errQueryParams
	}

	return net.JoinHostPort(host, port), nil
}

func connectTarget(addr string, ch chan net.Conn) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		ch <- nil
		return
	}
	ch <- conn
}