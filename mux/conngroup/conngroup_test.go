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

func (t *testConn) Read(p []byte) (n int, err error) {
	return copy(p, <-t.dataChan), nil
}

func (t *testConn) Write(p []byte) (n int, err error) {
	buf := make([]byte, len(p))
	n = copy(buf, p)
	t.dataChan <- buf
	return n, nil
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
		buf := make([]byte, 1024)
		for i := 0; i < 10; i++ {
			n, _ := cg.ReadMessage(buf)
			res[string(buf[:n])] = true
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
