package socks5

import (
	"errors"
	"log"
	"net"
	"strconv"
	"sync"
	"util"
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
		// Fake bind addr.
		0x00, 0x00, 0x00, 0x00,
		// Fake bind port.
		0x10, 0x10,
	}

	// buffer pool.
	bufPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, 257)
		},
	}

	// errors.
	errInvalid   = errors.New("Invalid target info")
	errNotSocks5 = errors.New("Not Socks5")
)

// Target target server info.
type Target struct {
	Host string
	Port string
}

// ConnHandler handle the connection
// after socks5 handshake succeeded.
type ConnHandler = func(net.Conn, *Target)

// ListenAndServe create a socks5 server.
func ListenAndServe(socks5Addr string, handler ConnHandler) {
	l, err := net.Listen("tcp", socks5Addr)
	if err != nil {
		log.Fatalln("Socks5 服务器监听失败，地址有误或端口被占用？")
	}

	log.Println("Socks5 server listening at", socks5Addr)

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
	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)

	// consult.
	_, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	if buf[0] != 0x05 {
		return nil, errNotSocks5
	}

	// reply for consult.
	_, err = conn.Write(consultRep)
	if err != nil {
		return nil, err
	}

	// get target info.
	n, err := conn.Read(buf)
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
	port := int(buf[l-2]) << 8 | int(buf[l-1])
	target.Port = strconv.Itoa(port)

	switch buf[0] {
	case ipv4:
		if l < 7 {
			return nil, errInvalid
		}
		target.Host = net.IP(buf[1:5]).String()

	case domain:
		domainLen := int(buf[1])
		if l < 4+domainLen {
			return nil, errInvalid
		}
		host := string(buf[2 : domainLen+2])
		ok, err := util.IsDomain(host)
		if !ok || err != nil {
			return nil, errInvalid
		}
		target.Host = host

	case ipv6:
		if l < 19 {
			return nil, errInvalid
		}
		target.Host = net.IP(buf[1:17]).String()

	default:
		return nil, errInvalid
	}

	return target, nil
}
