package main

import (
	"client"
	"flag"
	"log"
)

func main() {
	socks5Addr := flag.String("l", "127.0.0.1:6666", "Local socks5 server address.")
	serverAddr := flag.String("r", "127.0.0.1:80", "Falcon-WS server address.")
	fakeHost := flag.String("fh", "", "Fake 'Host' field for request header.")
	secure := flag.Bool("s", false, "Secure flag for enable https.")
	flag.Parse()

	config := &client.Config{
		Socks5Addr: *socks5Addr,
		ServerAddr: *serverAddr,
		FakeHost: *fakeHost,
		Secure: *secure,
	}

	log.Println("Socks5 server listening at", *socks5Addr)
	client.NewClient(config)
}
