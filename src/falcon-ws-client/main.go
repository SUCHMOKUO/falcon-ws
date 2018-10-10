package main

import (
	"client"
	"flag"
	"log"
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
		"Flag for enable dns cache. if it's true, it will lookup the ip of host and cache it.")

	flag.Parse()

	config := &client.Config{
		Socks5Addr: *socks5Addr,
		ServerAddr: *serverAddr,
		FakeHost: *fakeHost,
		Secure: *secure,
		Lookup: *lookup,
	}

	log.Println("Socks5 server listening at", *socks5Addr)
	client.NewClient(config)
}
