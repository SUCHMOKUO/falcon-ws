package main

import (
	"flag"
	"log"

	"github.com/SUCHMOKUO/falcon-ws/server"
)

func main() {
	addr := flag.String("l", ":80", "Listen address.")
	flag.Parse()
	log.Println("Server listening at", *addr)
	server.ListenAndServe(*addr)
}
