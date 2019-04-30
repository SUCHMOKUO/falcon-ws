package util

import (
	"encoding/base64"
	"errors"
	"io"
	"net"
	"regexp"
	"sync"
)

// Encode string using base64.
func Encode(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

// Decode base64 string.
func Decode(str string) (string, error) {
	bytes, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// IsDomain detect if value match the format of domain.
func IsDomain(host string) bool {
	ok, _ := regexp.MatchString(`\.[a-z]{2,}$`, host)
	return ok
}

// IsIPv4 detect if ip is ipv4.
func IsIPv4(ip net.IP) bool {
	return ip != nil && ip.To4() != nil
}

// IsIPv6 detect if ip is ipv6.
func IsIPv6(ip net.IP) bool {
	return ip != nil && ip.To4() == nil
}

type detector = func(net.IP) bool

// findIP returns the first ip matched detector.
func findIP(ips []net.IP, d detector) net.IP {
	for _, ip := range ips {
		if d(ip) {
			return ip
		}
	}
	return nil
}

var errNoIPv4 = errors.New("proxy server only has ipv6 address! please enable '-6' flag")

// Lookup return ip address of host.
func Lookup(host string, ipv6 bool) (net.IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	if ipv6 {
		// use ipv6 first.
		if ip := findIP(ips, IsIPv6); ip != nil {
			return ip, nil
		}
	}

	ip := findIP(ips, IsIPv4)
	if ip == nil {
		return nil, errNoIPv4
	}
	return ip, err
}

// buffer pool.
var bufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 1500)
	},
}

func Copy(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()

	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)

	for {
		n, err := src.Read(buf)
		if n > 0 {
			_, err := dst.Write(buf[:n])
			if err != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}
