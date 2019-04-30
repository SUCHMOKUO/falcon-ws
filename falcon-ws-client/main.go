package main

import (
	"flag"

	"github.com/SUCHMOKUO/falcon-ws/client"
)

func main() {
	var (
		socks5Addr string
		serverAddr string
		fakeHost string
		userAgent string
		secure bool
		lookup bool
		ipv6 bool
	)

	flag.StringVar(&socks5Addr,
		"l",
		"127.0.0.1:6666",
		"Local socks5 server address.")

	flag.StringVar(&serverAddr,
		"r",
		"127.0.0.1:80",
		"Falcon-WS server address.")

	flag.StringVar(&fakeHost,
		"fh",
		"",
		"Fake 'Host' field for request header.")

	flag.StringVar(&userAgent,
		"ua",
		"",
		"Fake 'User-Agent' field for request header.")

	flag.BoolVar(&secure,
		"secure",
		false,
		"Secure flag for enable https.")

	flag.BoolVar(&lookup,
		"lookup",
		true,
		"Flag for enable dns cache. if sets to 'true', it will lookup the server ip of host and cache it.")

	flag.BoolVar(&ipv6,
		"6",
		false,
		"Flag for enable ipv6. if sets to 'true', it will use ipv6 address (if it has) of proxy server first.")

	flag.Parse()

	client.New(&client.Config{
		Socks5Addr: socks5Addr,
		ServerAddr: serverAddr,
		FakeHost: fakeHost,
		UserAgent: userAgent,
		Secure: secure,
		Lookup: lookup,
		IPv6: ipv6,
	}).ListenAndServe()
}
