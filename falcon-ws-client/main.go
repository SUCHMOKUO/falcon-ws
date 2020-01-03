package main

import (
	"flag"

	"github.com/SUCHMOKUO/falcon-ws/client"
)

func main() {
	var (
		socks5Addr    string
		serverHost    string
		serverPort    string
		lookup        bool
		ipv6          bool
		password      string
		connGroupSize int
	)

	flag.StringVar(&socks5Addr,
		"socks5",
		"127.0.0.1:6666",
		"Local socks5 server address.")

	flag.StringVar(&serverHost,
		"host",
		"localhost",
		"Falcon-WS server host.")

	flag.StringVar(&serverPort,
		"port",
		"443",
		"Falcon-WS server port.")

	flag.BoolVar(&lookup,
		"lookup",
		true,
		"Flag for enable dns cache. if sets to 'true', it will lookup the server ip of host and cache it.")

	flag.BoolVar(&ipv6,
		"6",
		false,
		"Flag for enable ipv6. if sets to 'true', it will use ipv6 address (if it has) of proxy server first.")

	flag.StringVar(&password,
		"password",
		"password",
		"Password of proxy service.")

	flag.IntVar(&connGroupSize,
		"conngroup-size",
		5,
		"Size of connection group.")

	flag.Parse()

	client.New(&client.Config{
		Socks5Addr:    socks5Addr,
		ServerHost:    serverHost,
		ServerPort:    serverPort,
		Lookup:        lookup,
		IPv6:          ipv6,
		Password:      password,
		ConnGroupSize: connGroupSize,
	}).ListenAndServe()
}
