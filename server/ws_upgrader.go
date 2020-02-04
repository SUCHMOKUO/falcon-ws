package server

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/SUCHMOKUO/falcon-ws/messageconn"
	"github.com/gorilla/websocket"
)

const (
	wsTimeout = 10 * time.Second
	wsBufSize = 10300
)

var (
	// websocket upgrader.
	upgrader = websocket.Upgrader{
		HandshakeTimeout: wsTimeout,
		ReadBufferSize:   wsBufSize,
		WriteBufferSize:  wsBufSize,
		WriteBufferPool:  &sync.Pool{
			New: func() interface{} {
				return make([]byte, wsBufSize)
			},
		},
	}

	errUpgradeFail = errors.New("ws upgrade fail")
)

func wsUpgrade(ctx *Ctx) (err error, code int) {
	ws, err := upgrader.Upgrade(ctx.w, ctx.r, nil)
	if err != nil {
		return errUpgradeFail, http.StatusBadRequest
	}
	ctx.Data["conn"] = messageconn.NewWSConn(ws)
	return
}
