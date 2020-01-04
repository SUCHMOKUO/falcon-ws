package messageconn

import "github.com/gorilla/websocket"

// Conn is a message-based connection interface.
type Conn interface {
	ReadMessage() (msg []byte, err error)
	WriteMessage(msg []byte) error
	Close() error
}

type WSConn struct {
	conn *websocket.Conn
}

func NewWSConn(conn *websocket.Conn) Conn {
	return &WSConn{conn}
}

func (t *WSConn) ReadMessage() ([]byte, error) {
	_, data, err := t.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t *WSConn) WriteMessage(p []byte) error {
	return t.conn.WriteMessage(websocket.BinaryMessage, p)
}

func (t *WSConn) Close() error {
	return t.conn.Close()
}
