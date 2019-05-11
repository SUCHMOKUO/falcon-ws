package util

import "encoding/base64"

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
