package util

import (
	"net"
	"regexp"
)

// IsDomain detect if value match the format of domain.
func IsDomain(host string) (bool, error) {
	return regexp.MatchString(`\.[a-z]{2,}$`, host)
}

// IsIPv4 detect if ip is ipv4.
func IsIPv4(ip net.IP) bool {
	return ip.To4() != nil
}

// IsIPv6 detect if ip is ipv6.
func IsIPv6(ip net.IP) bool {
	return ip.To4() == nil
}

type detector = func(net.IP) bool

// FindIP returns the first ip matched detector.
func FindIP(ips []net.IP, d detector) net.IP {
	for _, ip := range ips {
		if d(ip) {
			return ip
		}
	}
	return nil
}