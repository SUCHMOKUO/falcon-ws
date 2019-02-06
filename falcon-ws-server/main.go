package main

import (
	"flag"
	"log"

	"github.com/SUCHMOKUO/falcon-ws/server"
)

func main() {
	port := flag.String("p", "80", "Listen port.")
	flag.Parse()
	log.Println("Server listening at", *port)
	server.NewServer(*port)
}
