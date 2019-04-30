package tunnel

import "github.com/gorilla/websocket"

const dataT = websocket.BinaryMessage

// Tunnel is the abstract of websocket connection.
type Tunnel struct {
	websocket.Conn
}

// Read implements io.Reader.
// length of p must be greater than 1500.
func (t *Tunnel) Read(p []byte) (int, error) {
	_, data, err := t.ReadMessage()
	if err != nil {
		return 0, err
	}
	return copy(p, data), nil
}

// Write implements io.Writer.
func (t *Tunnel) Write(p []byte) (int, error) {
	err := t.WriteMessage(dataT, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}