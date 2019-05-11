package server

import (
	"log"
	"net/http"
)

// ListenAndServe create a falcon-ws server.
func ListenAndServe(addr string) {
	http.HandleFunc("/free", handleProxyReq)
	http.HandleFunc("/location", handleLocationReq)
	log.Fatal(http.ListenAndServe(addr, nil))
}
