package util

import (
	"errors"
	"github.com/gorilla/websocket"
	"net"
	"regexp"
	"sync"
)

// IsDomain detect if value match the format of domain.
func IsDomain(host string) bool {
	ok, _ := regexp.MatchString(`\.[a-z]{2,}$`, host)
	return ok
}

// IsIPv4 detect if ip is ipv4.
func IsIPv4(ip net.IP) bool {
	return ip != nil && ip.To4() != nil
}

// IsIPv6 detect if ip is ipv6.
func IsIPv6(ip net.IP) bool {
	return ip != nil && ip.To4() == nil
}

type detector = func(net.IP) bool

// findIP returns the first ip matched detector.
func findIP(ips []net.IP, d detector) net.IP {
	for _, ip := range ips {
		if d(ip) {
			return ip
		}
	}
	return nil
}

var errNoIPv4 = errors.New("proxy server only has ipv6 address! please enable '-6' flag")

// Lookup return ip address of host.
func Lookup(host string, ipv6 bool) (net.IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	if ipv6 {
		// use ipv6 first.
		if ip := findIP(ips, IsIPv6); ip != nil {
			return ip, nil
		}
	}

	ip := findIP(ips, IsIPv4)
	if ip == nil {
		return nil, errNoIPv4
	}
	return ip, err
}

// buffer pool.
var bufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func WSToConn(conn net.Conn, ws *websocket.Conn) {
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
			err = ws.WriteMessage(websocket.BinaryMessage, buf[:n])
			if err != nil {
				return
			}
		}
	}
}

func ConnToWS(ws *websocket.Conn, conn net.Conn) {
	defer ws.Close()
	defer conn.Close()

	for {
		msgT, data, err := ws.ReadMessage()

		if msgT != websocket.BinaryMessage || err != nil {
			return
		}

		_, err = conn.Write(data)

		if err != nil {
			return
		}
	}
}