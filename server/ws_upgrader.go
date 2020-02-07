package server

import (
	"errors"
	"github.com/SUCHMOKUO/falcon-ws/configs"
	"github.com/SUCHMOKUO/falcon-ws/messageconn"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var (
	// websocket upgrader.
	upgrader = websocket.Upgrader{
		HandshakeTimeout: configs.Timeout,
		ReadBufferSize:   configs.MaxPackageSize,
		WriteBufferSize:  configs.MaxPackageSize,
		WriteBufferPool:  &sync.Pool{
			New: func() interface{} {
				return make([]byte, configs.MaxPackageSize)
			},
		},
	}
)

var (
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
