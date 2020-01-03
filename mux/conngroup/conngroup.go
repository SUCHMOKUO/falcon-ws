package conngroup

import (
	"context"
	"github.com/SUCHMOKUO/falcon-ws/messageconn"
	"github.com/SUCHMOKUO/falcon-ws/util"
	"io"
)

const (
	chanBufSize = 128
)

// ConnGroup is a group of few connections.
type ConnGroup struct {
	// connList is the list of ws connections.
	connList *ConnList

	readChan  chan []byte
	writeChan chan []byte

	closed bool

	// context for manage goroutines.
	ctx    context.Context
	cancel context.CancelFunc
}

func (cg *ConnGroup) Size() int {
	return cg.connList.Size()
}

func (cg *ConnGroup) ReadMessage() ([]byte, error) {
	data, ok := <-cg.readChan
	if !ok {
		// connection group closed.
		return nil, io.EOF
	}
	return data, nil
}

func (cg *ConnGroup) WriteMessage(p []byte) error {
	cg.writeChan <- util.CopyBuf(p)
	return nil
}

func (cg *ConnGroup) Close() {
	cg.closed = true
	cg.cancel()
	for i := cg.Size(); i > 0; i-- {
		cg.connList.Next().Close()
	}
	close(cg.readChan)
	close(cg.writeChan)
}

func (cg *ConnGroup) AddConn(netConn messageconn.Conn) {
	conn := NewConn(netConn)
	cg.connList.Add(conn)
	go cg.handleRead(conn)
}

func (cg *ConnGroup) handleWrite() {
	for {
		select {
		case <-cg.ctx.Done():
			return
		case msg, ok := <-cg.writeChan:
			if !ok {
				return
			}
			for !cg.closed {
				conn := cg.connList.Next()
				if conn.Usable() {
					conn.WriteMessage(msg)
					break
				}
			}
		}
	}
}

func (cg *ConnGroup) enableHeartBeat() {
	cg.connList.ForEach(func(conn *Conn) {
		// TODO add heart beat feature for conn.
	})
}
