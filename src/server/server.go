package server

import (
	"log"
	"net"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
	"encoding/base64"
	"sync"
)

const (
	// transmission type
	dataT   = websocket.BinaryMessage
	// connection time out
	timeout = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: timeout,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	WriteBufferPool:  &sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	},
}

// NewServer create a falcon-ws server.
func NewServer(port string) {
	http.HandleFunc("/free", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query()
	hostEnc := target.Get("h")
	portEnc := target.Get("p")

	if hostEnc== "" || portEnc == "" {
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

	// log.Println(host, port)

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

	go recive(conn, ws)
	go send(ws, conn)
}

func connectTarget(addr string, ch chan net.Conn) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		ch <- nil
		return
	}
	ch <- conn
}

func recive(conn net.Conn, ws *websocket.Conn) {
	defer ws.Close()
	defer conn.Close()

	buf := make([]byte, 1024)
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

func send(ws *websocket.Conn, conn net.Conn) {
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