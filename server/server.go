package server

import (
	"encoding/base64"
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
	bufSize = 1024
)

var (
	// websocket upgrader.
	upgrader = websocket.Upgrader{
		HandshakeTimeout: timeout,
		ReadBufferSize:   bufSize,
		WriteBufferSize:  bufSize,
		WriteBufferPool:  new(sync.Pool),
	}
)

// NewServer create a falcon-ws server.
func NewServer(port string) {
	http.HandleFunc("/free", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query()
	hostEnc := target.Get("h")
	portEnc := target.Get("p")

	if hostEnc == "" || portEnc == "" {
		http.NotFound(w, r)
		return
	}

	// url decode.
	hostS, err := base64.URLEncoding.DecodeString(hostEnc)
	portS, err := base64.URLEncoding.DecodeString(portEnc)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	host := string(hostS)
	port := string(portS)

	// log.Println(host, port, r.Host)

	addr := net.JoinHostPort(host, port)
	ch := make(chan net.Conn, 1)
	go connectTarget(addr, ch)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conn := <-ch
	if conn == nil {
		ws.Close()
		return
	}

	go util.WSToConn(conn, ws)
	go util.ConnToWS(ws, conn)
}

func connectTarget(addr string, ch chan net.Conn) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		ch <- nil
		return
	}
	ch <- conn
}
