package socks5

import (
	"errors"
	"log"
	"net"
	"strconv"
	"regexp"
)

const (
	ipv4   = 0x01
	domain = 0x03
	ipv6   = 0x04
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
		// Fake bind addr
		0x00, 0x00, 0x00, 0x00,
		// Fake bind port
		0x10, 0x10,
	}
)

// Target target server info.
type Target struct {
	Host string
	Port string
}

// ConnHandler connection handler function
// after socks5 handshake succeeded.
type ConnHandler = func(net.Conn, *Target)

// ListenAndServe create a socks5 server.
func ListenAndServe(socks5Addr string, handler ConnHandler) {
	l, err := net.Listen("tcp", socks5Addr)
	if err != nil {
		log.Fatalln("Socks5 服务器监听失败，地址有误或端口被占用？")
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go func() {
			target, err := socks5Handshake(conn)
			if err != nil {
				conn.Close()
				log.Println("Socks5 handshake error:", err)
				return
			}
			// start proxy.
			handler(conn, target)
		}()
	}
}

func socks5Handshake(conn net.Conn) (*Target, error) {
	buf := make([]byte, 257)
	// consult.
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	if buf[0] != 0x05 {
		return nil, errors.New("Not Socks5")
	}

	// reply for consult.
	_, err = conn.Write(consultRep)
	if err != nil {
		return nil, err
	}

	// get target info.
	n, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// reply for info.
	_, err = conn.Write(infoRep)
	if err != nil {
		return nil, err
	}

	return parseTargetInfo(buf[3:n])
}

func parseTargetInfo(buf []byte) (*Target, error) {
	target := new(Target)

	l := len(buf)
	if l < 7 {
		return nil, errors.New("Invalid target info")
	}

	port := int(buf[l-2])*256 + int(buf[l-1])
	target.Port = strconv.Itoa(port)

	switch buf[0] {
	case ipv4:
		target.Host = net.IP(buf[1:5]).String()

	case domain:
		domainLen := buf[1]
		host := string(buf[2 : domainLen+2])
		ok, err := isDomain(host)
		if !ok || err != nil {
			return nil, errors.New("Invalid target info")
		}
		target.Host = host

	case ipv6:
		target.Host = net.IP(buf[1:17]).String()

	default:
		return nil, errors.New("Invalid target info")
	}

	return target, nil
}

// detect if value match the format of domain.
func isDomain(host string) (bool, error) {
	return regexp.MatchString(`\.[a-z]{2,}$`, host)
}