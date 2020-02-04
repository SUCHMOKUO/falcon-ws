package socks5

import (
	"errors"
	"fmt"
	"github.com/SUCHMOKUO/falcon-ws/util"
	"log"
	"net"
)

const (
	IPV4   = 0x01
	DOMAIN = 0x03
	IPV6   = 0x04
)

var (
	consultRep = []byte{0x05, 0x00}

	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	infoRep = []byte{
		0x05, 0x00, 0x00, 0x01,
		// Fake bind addr.
		0x00, 0x00, 0x00, 0x00,
		// Fake bind port.
		0x10, 0x10,
	}
)

var (
	// errors.
	errInvalid   = errors.New("invalid target info")
	errNotSocks5 = errors.New("not socks5")
)

// ProxyFunc handle the connection
// after socks5 handshake succeeded.
type ProxyFunc = func(socksConn net.Conn, targetAddr string)

// ListenAndServe create a socks5 server.
func ListenAndServe(socks5Addr string, p ProxyFunc) {
	l, err := net.Listen("tcp", socks5Addr)
	if err != nil {
		log.Fatalln("[Socks5] 服务器监听失败，地址有误或端口被占用？")
	}

	log.Println("[Socks5] server listening at", socks5Addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go handleSocks5(conn, p)
	}
}

func handleSocks5(socksConn net.Conn, p ProxyFunc) {
	targetAddr, err := socks5Handshake(socksConn)
	if err != nil {
		socksConn.Close()
		log.Println("[Socks5] handshake error:", err)
		return
	}
	// start proxy.
	p(socksConn, targetAddr)
}

func socks5Handshake(conn net.Conn) (string, error) {
	var bufArr [257]byte
	buf := bufArr[:]

	// consult.
	_, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	if buf[0] != 0x05 {
		return "", errNotSocks5
	}

	// reply for consult.
	_, err = conn.Write(consultRep)
	if err != nil {
		return "", err
	}

	// get target info.
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// reply for info.
	_, err = conn.Write(infoRep)
	if err != nil {
		return "", err
	}

	// parse target info.
	info := buf[3:n]

	l := len(info)
	port := int(info[l-2])<<8 | int(info[l-1])
	var hostStr string

	switch info[0] {
	case IPV4:
		if l < 7 {
			return "", errInvalid
		}
		hostStr = net.IP(info[1:5]).String()

	case DOMAIN:
		domainLen := int(info[1])
		if l < 4+domainLen {
			return "", errInvalid
		}
		host := string(info[2 : domainLen+2])
		if !util.IsValidHost(host) {
			return "", errInvalid
		}
		hostStr = host

	case IPV6:
		if l < 19 {
			return "", errInvalid
		}
		hostStr = net.IP(info[1:17]).String()

	default:
		return "", errInvalid
	}

	return fmt.Sprintf("%s:%d", hostStr, port), nil
}
