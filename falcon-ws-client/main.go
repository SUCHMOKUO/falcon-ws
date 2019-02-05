package main

import (
	"client"
	"flag"
)

func main() {
	socks5Addr := flag.String(
		"l",
		"127.0.0.1:6666",
		"Local socks5 server address.")

	serverAddr := flag.String(
		"r",
		"127.0.0.1:80",
		"Falcon-WS server address.")

	fakeHost := flag.String(
		"fh",
		"",
		"Fake 'Host' field for request header.")

	secure := flag.Bool(
		"secure",
		false,
		"Secure flag for enable https.")

	lookup := flag.Bool(
		"lookup",
		true,
		"Flag for enable dns cache. if sets to 'true', it will lookup the server ip of host and cache it.")

	ipv6 := flag.Bool(
		"6",
		false,
		"Flag for enable ipv6. if sets to 'true', it will use ipv6 address (if it has) of proxy server first.")

	flag.Parse()

	config := &client.Config{
		Socks5Addr: *socks5Addr,
		ServerAddr: *serverAddr,
		FakeHost:   *fakeHost,
		Secure:     *secure,
		Lookup:     *lookup,
		IPv6:       *ipv6,
	}

	client.NewClient(config)
}
