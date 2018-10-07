package main

import (
	"client"
	"flag"
	"log"
)

func main() {
	socks5Addr := flag.String("l", "127.0.0.1:6666", "Local socks5 server address.")
	serverAddr := flag.String("r", "127.0.0.1:80", "Falcon-WS server address.")
	flag.Parse()
	log.Println("Socks5 server listening at", *socks5Addr)
	client.NewClient(*socks5Addr, *serverAddr)
}
