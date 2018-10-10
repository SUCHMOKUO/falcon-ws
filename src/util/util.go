package util

import "regexp"

// IsDomain detect if value match the format of domain.
func IsDomain(host string) (bool, error) {
	return regexp.MatchString(`\.[a-z]{2,}$`, host)
}
