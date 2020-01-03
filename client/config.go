package client

// Config of client.
type Config struct {
	// local socks5 address for listening.
	Socks5Addr string

	// falcon server host.
	ServerHost string

	// falcon server port.
	ServerPort string

	// lookup flag. if it's false, client will
	// not lookup the server ip and cache it.
	Lookup bool

	// ipv6 flag. if it's true, the ipv6 address
	// of proxy server will be used first.
	IPv6 bool

	// password for proxy service.
	Password string

	// size of connection group.
	ConnGroupSize int
}
