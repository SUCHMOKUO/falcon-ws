package server

var globalConfig *Config

type Config struct {
	// server listen address.
	Addr string

	// certification file path.
	Cert string

	// private key file path.
	PrivateKey string

	// key for digital signature.
	SignatureKey string

	// password for proxy service.
	Password string
}
