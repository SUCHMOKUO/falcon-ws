package conngroup

import (
	"testing"
)

type testIO struct {
	id int
}

func (t testIO) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (t testIO) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (t testIO) Close() error {
	return nil
}

func TestConnList_Add(t *testing.T) {
	cl := NewList()
	cl.Add(NewConn(nil))
	cl.Add(NewConn(nil))
	cl.Add(NewConn(nil))
	if cl.Size() != 3 {
		t.Error("size should be 3")
	}
}

func TestConnList_GetCurConn(t *testing.T) {
	cl := NewList()
	cl.Add(NewConn(testIO{0}))
	cl.Add(NewConn(testIO{1}))
	cl.Add(NewConn(testIO{2}))
	if cl.GetCurConn().conn.(testIO).id != 2 {
		t.Error("id should be 2")
	}
}

func TestConnList_Next(t *testing.T) {
	cl := NewList()
	cl.Add(NewConn(testIO{0}))
	cl.Add(NewConn(testIO{1}))
	cl.Add(NewConn(testIO{2}))
	cl.Next()
	if cl.GetCurConn().conn.(testIO).id != 1 {
		t.Error("id should be 1")
	}
}

func TestConnList_Remove(t *testing.T) {
	cl := NewList()
	cl.Add(NewConn(testIO{1}))
	c := NewConn(testIO{2})
	cl.Add(c)
	cl.Remove(c)
	if cl.Size() != 1 {
		t.Error("size should be 1 after remove")
	}
	cl.ForEach(func(conn *Conn) {
		if conn.conn.(testIO).id == 2 {
			t.Error("remove fail")
		}
	})

	cl = NewList()
	c = NewConn(testIO{0})
	cl.Add(c)
	cl.Remove(c)
	if cl.Size() != 0 {
		t.Error("size should be 0 after remove")
	}
	cl.ForEach(func(conn *Conn) {
		t.Error("should not be executed")
	})

	cl = NewList()
	cl.Add(NewConn(testIO{0}))
	cl.Remove(NewConn(testIO{1}))
	if cl.Size() != 1 {
		t.Error("should not remove anything")
	}
}