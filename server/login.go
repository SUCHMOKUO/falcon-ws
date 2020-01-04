package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/SUCHMOKUO/falcon-ws/util"
)

var (
	errNoPassword = errors.New("no password")
	errWrongPassword = errors.New("wrong password")
)

func login(ctx *Ctx) (err error, code int) {
	p := ctx.r.PostFormValue("password")
	if p == "" {
		return errNoPassword, http.StatusBadRequest
	}
	if p != globalConfig.Password {
		return errWrongPassword, http.StatusForbidden
	}
	ctx.w.Write([]byte(newToken()))
	return
}

func newToken() string {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": util.RandomUint32String(),
	}).SignedString([]byte(globalConfig.SignatureKey))
	if err != nil {
		log.Fatalln("[JWT Signature]", err)
	}
	return token
}
