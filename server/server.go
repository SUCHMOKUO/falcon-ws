package server

import (
	"log"
	"net/http"
)

// ListenAndServe create a falcon-ws server.
func ListenAndServe(config *Config) {
	globalConfig = config

	loginHandler := new(HttpHandler)
	loginHandler.Use(login)
	http.HandleFunc("/login", loginHandler.GetHandleFunc())

	proxyHandler := new(HttpHandler)
	proxyHandler.Use(auth)
	proxyHandler.Use(wsUpgrade)
	proxyHandler.Use(sessionManager)
	http.HandleFunc("/free", proxyHandler.GetHandleFunc())

	log.Fatalln(http.ListenAndServeTLS(config.Addr, config.Cert, config.PrivateKey, nil))
}
