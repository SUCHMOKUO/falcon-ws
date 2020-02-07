package server

import (
	"github.com/SUCHMOKUO/falcon-ws/configs"
	"log"
	"net"

	"github.com/SUCHMOKUO/falcon-ws/mux/stream"
	"github.com/SUCHMOKUO/falcon-ws/util"
)

func proxy(s *stream.Stream) {
	buf := make([]byte, 260)
	n, err := s.Read(buf)
	if err != nil {
		s.Close()
		return
	}

	addr := string(buf[:n])

	conn, err := net.DialTimeout("tcp", addr, configs.Timeout)
	if err != nil {
		s.Close()
		return
	}

	log.Println("[Proxy] start proxy target:", addr)

	go util.CopyIO(conn, s)
	go util.CopyIO(s, conn)
}
