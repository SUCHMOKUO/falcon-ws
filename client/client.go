package client

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/SUCHMOKUO/falcon-ws/socks5"
	"github.com/SUCHMOKUO/falcon-ws/tunnel"
	"github.com/SUCHMOKUO/falcon-ws/util"
	"github.com/gorilla/websocket"
)

const (
	// connection time out
	timeout = 10 * time.Second

	// buffer size.
	bufSize = 1500
)

var (
	// fake User-Agent.
	defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML"
)

// Config of client.
type Config struct {
	// local socks5 address for listening.
	Socks5Addr string

	// falcon server address. (host:port)
	ServerAddr string

	// fake 'Host' field in request header
	// for against qos.
	FakeHost string

	// fake User-Agent.
	UserAgent string

	// secure flag. true for using https
	// instead of http.
	Secure bool

	// lookup flag. if it's false, client will
	// not lookup the server ip and cache it.
	Lookup bool

	// ipv6 flag. if it's true, the ipv6 address
	// of proxy server will be used first.
	IPv6 bool
}

// Client
type Client struct {
	// config
	config *Config

	// extra server info.
	host   string
	port   string
	wsAddr string

	// http header
	header http.Header

	// websocket dialer.
	dialer websocket.Dialer
}

func (c *Client) completeHostAndPort() {
	host, port, err := net.SplitHostPort(c.config.ServerAddr)
	if err != nil {
		log.Fatalln("server address format error!")
	}
	c.host = host
	c.port = port
}

func (c *Client) completeWSAddr() {
	// prepend schema.
	var schema string
	if c.config.Secure {
		schema = "wss://"
	} else {
		schema = "ws://"
	}

	// lookup server ip.
	if !util.IsDomain(c.host) || !c.config.Lookup {
		c.wsAddr = schema + c.config.ServerAddr
		return
	}
	log.Println("lookup for server", c.host)
	ip, err := util.Lookup(c.host, c.config.IPv6)
	if err != nil {
		log.Fatalln(err)
	}
	var ipStr string
	if util.IsIPv4(ip) {
		ipStr = ip.String()
	}
	if util.IsIPv6(ip) {
		ipStr = "[" + ip.String() + "]"
	}
	c.wsAddr = schema + ipStr + ":" + c.port
}

// init tls setting.
func (c *Client) completeSecureSetting() {
	if !c.config.Secure {
		return
	}
	tlsCfg := new(tls.Config)
	tlsCfg.ServerName = c.host
	c.dialer.TLSClientConfig = tlsCfg
}

// init fake header.
func (c *Client) completeHeader() {
	reqHeader := http.Header{}
	if c.config.FakeHost != "" && !c.config.Secure {
		// add fake host field.
		reqHeader.Set("Host", c.config.FakeHost)
	} else {
		reqHeader.Set("Host", c.host)
	}
	// fake user-agent field.
	if c.config.UserAgent != "" {
		reqHeader.Set("User-Agent", c.config.UserAgent)
	} else {
		reqHeader.Set("User-Agent", defaultUserAgent)
	}
	c.header = reqHeader
}

// New return new falcon client instance.
func New(cfg *Config) *Client {
	c := new(Client)
	c.config = cfg
	c.dialer = websocket.Dialer{
		HandshakeTimeout: timeout,
		ReadBufferSize:   bufSize,
		WriteBufferSize:  bufSize,
		WriteBufferPool:  new(sync.Pool),
	}
	c.completeHostAndPort()
	c.completeWSAddr()
	c.completeSecureSetting()
	c.completeHeader()
	return c
}

// CreateTunnel create a tunnel through falcon server to target.
func (c *Client) CreateTunnel(targetHost, targetPort string) (io.ReadWriteCloser, error) {
	// url encode.
	host := util.Encode(targetHost)
	port := util.Encode(targetPort)
	url := c.wsAddr + "/free?h=" + host + "&p=" + port
	ws, res, err := c.dialer.Dial(url, c.header)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 404 {
		return nil, errors.New("Empty target address or port.")
	}
	t := &tunnel.Tunnel{*ws}
	return t, nil
}

// ListenAndServe create a local socks5 server.
func (c *Client) ListenAndServe() {
	log.Println("falcon server:", c.wsAddr)
	log.Println("use host:", c.header.Get("Host"))
	log.Println("use user-agent:", c.header.Get("User-Agent"))

	handleConn := func(conn net.Conn, target *socks5.Target) {
		t, err := c.CreateTunnel(target.Host, target.Port)
		if err != nil {
			log.Println("dial proxy server error:", err)
			return
		}
		go util.Copy(t, conn)
		go util.Copy(conn, t)
	}

	// start socks5 server.
	socks5.ListenAndServe(c.config.Socks5Addr, handleConn)
}
