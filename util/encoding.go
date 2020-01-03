package util

import "encoding/base64"

// EncodeBase64 string using base64.
func EncodeBase64(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

// DecodeBase64 base64 string.
func DecodeBase64(str string) (string, error) {
	bytes, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
