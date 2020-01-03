package client

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/SUCHMOKUO/falcon-ws/messageconn"
	"github.com/SUCHMOKUO/falcon-ws/mux"
	"github.com/SUCHMOKUO/falcon-ws/mux/conngroup"
	"github.com/SUCHMOKUO/falcon-ws/socks5"
	"github.com/SUCHMOKUO/falcon-ws/util"
	"github.com/gorilla/websocket"
)

const (
	// connection time out
	timeout = 10 * time.Second

	bufSize = 4096
)

// Client
type Client struct {
	// config
	config *Config

	// extra server info.
	wsUrl    string
	httpsUrl string

	// http header for create ws connection.
	header http.Header

	// websocket dialer.
	dialer websocket.Dialer

	m *mux.Mux
}

// NewWSConn return new falcon client instance.
func New(cfg *Config) *Client {
	c := new(Client)
	c.config = cfg
	c.dialer = websocket.Dialer{
		HandshakeTimeout: timeout,
		ReadBufferSize:   bufSize,
		WriteBufferSize:  bufSize,
		WriteBufferPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, bufSize)
			},
		},
		TLSClientConfig: &tls.Config{
			ServerName: cfg.ServerHost,
		},
	}
	c.header = http.Header{}
	c.completeUrl()
	c.login()
	c.createConnGroup()
	return c
}

// CreateProxyConn create a proxy connection through falcon server to target.
func (c *Client) CreateProxyConn(targetHost, targetPort string) (io.ReadWriteCloser, error) {
	targetUrl := util.EncodeBase64(targetHost + ":" + targetPort)
	s, err := c.m.NewStream()
	if err != nil {
		log.Fatalln("[Client]", err)
	}
	_, err = s.Write([]byte(targetUrl))
	if err != nil {
		return nil, err
	}
	return s, nil
}

// ListenAndServe create a local socks5 server.
func (c *Client) ListenAndServe() {
	log.Println("[Client] falcon server:", c.wsUrl)

	handleConn := func(conn net.Conn, target *socks5.Target) {
		t, err := c.CreateProxyConn(target.Host, target.Port)
		if err != nil {
			log.Println("[Create Conn]", err)
			conn.Close()
			return
		}
		go util.CopyIO(t, conn)
		go util.CopyIO(conn, t)
	}

	// start socks5 server.
	socks5.ListenAndServe(c.config.Socks5Addr, handleConn)
}

func (c *Client) completeUrl() {
	c.httpsUrl = "https://" + c.config.ServerHost + ":" + c.config.ServerPort
	// lookup server ip.
	if !util.IsDomain(c.config.ServerHost) || !c.config.Lookup {
		c.wsUrl = "wss://" + c.config.ServerHost + ":" + c.config.ServerPort
		return
	}
	log.Println("[Client] lookup for server:", c.config.ServerHost)
	ip, err := util.Lookup(c.config.ServerHost, c.config.IPv6)
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
	c.wsUrl = "wss://" + ipStr + ":" + c.config.ServerPort
}

func (c *Client) login() {
	form := url.Values{}
	form.Set("password", c.config.Password)
	resp, err := http.PostForm(c.httpsUrl+"/login", form)
	if err != nil {
		log.Fatalln("[Login]", err)
	}

	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("[Login]", err)
	}
	defer resp.Body.Close()

	tokenStr := string(token)
	c.header.Set("Authorization", "Bearer "+tokenStr)
	log.Println("[Login] Succeed")
}

func (c *Client) createConnGroup() {
	cg := conngroup.New(&conngroup.Config{
		Size: c.config.ConnGroupSize,
		GenConn: func() (messageconn.Conn, error) {
			ws, res, wsErr := c.dialer.Dial(c.wsUrl+"/free", c.header)
			if wsErr != nil {
				return nil, wsErr
			}
			code := res.StatusCode
			msg, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()
			log.Println("[Creating Connection Group]", code, string(msg))
			return messageconn.NewWSConn(ws), nil
		},
	})
	c.m = mux.New(cg)
}
