// +build server

package conngroup

import (
	"context"
)

// New return a new instance of ConnGroup.
func New() *ConnGroup {
	cg := new(ConnGroup)
	cg.connList = NewList()
	cg.readChan = make(chan []byte, chanBufSize)
	cg.writeChan = make(chan []byte, chanBufSize)
	cg.ctx, cg.cancel = context.WithCancel(context.Background())
	go cg.handleWrite()
	return cg
}

func (cg *ConnGroup) handleRead(conn *Conn) {
	for {
		select {
		case <-cg.ctx.Done():
			return
		default:
			msg, err := conn.ReadMessage()
			if err != nil {
				goto HandleConnClose
			}
			cg.readChan <- msg
		}
	}
HandleConnClose:
	cg.connList.Remove(conn)
}
