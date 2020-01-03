package server

import (
	"github.com/SUCHMOKUO/falcon-ws/messageconn"
	"sync"
)

var (
	sessions sync.Map
)

func sessionManager(ctx *Ctx) (err error, code int) {
	sessionId := ctx.Data["id"].(string)
	conn := ctx.Data["conn"].(messageconn.Conn)
	var s *session
	val, ok := sessions.Load(sessionId)
	if !ok {
		s = newSession(sessionId)
		sessions.Store(sessionId, s)
		go s.startProxyService()
		//log.Println("[Session Manager] new session:", sessionId)
	} else {
		s = val.(*session)
	}
	s.cg.AddConn(conn)
	//log.Println("[Session Manager] connection added to session", sessionId)
	return
}
