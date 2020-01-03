package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var (
	errUnexpectedSignMethod = errors.New("unexpected signing method")
	errInvalidToken = errors.New("invalid token")
	errNoAuth = errors.New("no authorization")
	errUnexpectedAuthFormat = errors.New("unexpected authorization format")
)

func auth(ctx *Ctx) (err error, code int) {
	a := ctx.r.Header.Get("Authorization")
	if a == "" {
		return errNoAuth, http.StatusUnauthorized
	}
	token, err := getToken(a)
	if err != nil {
		return err, http.StatusBadRequest
	}
	id, err := parseToken(token)
	if err != nil {
		return err, http.StatusForbidden
	}
	ctx.Data["id"] = id
	return
}

func getToken(authStr string) (token string, err error) {
	res := strings.Split(authStr, "Bearer ")
	if len(res) < 2 {
		return "", errUnexpectedAuthFormat
	}
	return res[1], nil
}

func parseToken(tokenStr string) (id string, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnexpectedSignMethod
		}
		return []byte(globalConfig.SignatureKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["id"].(string), nil
	} else {
		return "", errInvalidToken
	}
}