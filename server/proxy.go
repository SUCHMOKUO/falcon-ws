package server

import (
	"log"
	"net"
	"time"

	"github.com/SUCHMOKUO/falcon-ws/mux/stream"
	"github.com/SUCHMOKUO/falcon-ws/util"
)

const dialTimeout = 10 * time.Second

func proxy(s *stream.Stream) {
	buf := make([]byte, 512)
	n, err := s.Read(buf)
	if err != nil {
		s.Close()
		return
	}

	addr, err := util.DecodeBase64(string(buf[:n]))
	if err != nil {
		s.Close()
		return
	}

	conn, err := net.DialTimeout("tcp", addr, dialTimeout)
	if err != nil {
		s.Close()
		return
	}

	log.Println("[Proxy] start proxy target:", addr)

	go util.CopyIO(conn, s)
	go util.CopyIO(s, conn)
}
