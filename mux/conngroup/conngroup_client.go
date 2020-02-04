// +build !server

package conngroup

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/SUCHMOKUO/falcon-ws/messageconn"
)

var (
	globalConfig *Config
)

type Config struct {
	// Size is the number of connections of this group.
	Size int

	// GenConn is the connection generator function.
	GenConn func() (messageconn.Conn, error)

	// EnableHeartBeat for the ws connections.
	EnableHeartBeat bool

	// HeartBeatInterval is how often to send a heart beat message.
	HeartBeatInterval time.Duration
}

// NewWSConn return a new instance of ConnGroup.
func New(config *Config) *ConnGroup {
	globalConfig = config
	cg := new(ConnGroup)
	cg.connList = NewList()
	cg.readChan = make(chan []byte, chanBufSize)
	cg.writeChan = make(chan []byte, chanBufSize)
	cg.ctx, cg.cancel = context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(config.Size)
	for i := config.Size; i > 0; i-- {
		go func() {
			defer wg.Done()
			for {
				conn, err := config.GenConn()
				if err != nil {
					log.Println("[Creating Connection Group] Failed:", err)
					time.Sleep(time.Second * 5)
					continue
				}
				cg.AddConn(conn)
				log.Println("[Creating Connection Group] Succeed, current group size:", cg.Size())
				return
			}
		}()
	}
	wg.Wait()

	go cg.handleWrite()

	if config.EnableHeartBeat {
		cg.enableHeartBeat()
	}

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
	cg.handleBrokenConn(conn)
}

func (cg *ConnGroup) handleBrokenConn(conn *Conn) {
	if cg.closed {
		return
	}
	conn.SetUnusable()
	var newConn messageconn.Conn
	var err error
	for {
		newConn, err = globalConfig.GenConn()
		if err == nil {
			break
		}
	}
	conn.ReplaceConn(newConn)
	go cg.handleRead(conn)
	conn.SetUsable()
}
