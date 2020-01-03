package main

import (
	"flag"
	"log"

	"github.com/SUCHMOKUO/falcon-ws/server"
)

func main() {
	addr := flag.String("addr", ":443", "Listen address.")
	cert := flag.String("cert", "cert.pem", "Certificate file path.")
	privateKey := flag.String("privkey", "key.pem", "Private key file path.")
	signatureKey := flag.String("sigkey", "signatureKey", "SignatureKey for digital signature.")
	password := flag.String("password", "password", "Password for proxy service.")
	flag.Parse()

	config := server.Config{
		Addr:         *addr,
		Cert:         *cert,
		PrivateKey:   *privateKey,
		SignatureKey: *signatureKey,
		Password:     *password,
	}

	log.Println("[Server] listening at", *addr)
	server.ListenAndServe(&config)
}
