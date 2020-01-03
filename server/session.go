package server

import (
	"context"

	"github.com/SUCHMOKUO/falcon-ws/mux"
	"github.com/SUCHMOKUO/falcon-ws/mux/conngroup"
)

type session struct {
	id     string
	cg     *conngroup.ConnGroup
	mux    *mux.Mux
	ctx    context.Context
	cancel context.CancelFunc
}

func newSession(id string) *session {
	s := new(session)
	s.id = id
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.cg = conngroup.New()
	s.mux = mux.New(s.cg)
	return s
}

func (ss *session) startProxyService() {
	for {
		select {
		case <-ss.ctx.Done():
			return
		default:
			s, err := ss.mux.Accept()
			if err != nil {
				return
			}
			go proxy(s)
		}
	}
}
