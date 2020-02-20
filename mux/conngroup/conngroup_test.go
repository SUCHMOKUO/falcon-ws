// +build server

package conngroup

import (
	"strconv"
	"testing"
)

type testConn struct {
	dataChan chan []byte
}

func newTestConn() *testConn {
	return &testConn{
		dataChan: make(chan []byte, 8),
	}
}

func (t *testConn) ReadMessage() (msg []byte, err error) {
	return <-t.dataChan, nil
}

func (t *testConn) WriteMessage(msg []byte) error {
	t.dataChan <- msg
	return nil
}

func (t *testConn) Close() error {
	return nil
}

func TestConnGroup(t *testing.T) {
	cg := New()
	for i := 0; i < 5; i++ {
		cg.AddConn(newTestConn())
	}
	done := make(chan bool)
	res := make(map[string]bool, 10)
	go func() {
		for i := 0; i < 10; i++ {
			msg, _ := cg.ReadMessage()
			res[string(msg)] = true
		}
		done <- true
	}()
	for i := 0; i < 10; i++ {
		str := strconv.Itoa(i)
		cg.WriteMessage([]byte(str))
	}
	<-done
	for i := 0; i < 10; i++ {
		str := strconv.Itoa(i)
		if _, ok := res[str]; !ok {
			t.Error("result should contains:", str)
		}
	}
}
