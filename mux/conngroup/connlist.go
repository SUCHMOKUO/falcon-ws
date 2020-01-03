package conngroup

import (
	"container/ring"
	"sync"

	"github.com/SUCHMOKUO/falcon-ws/messageconn"
)

type Conn struct {
	lock   sync.Mutex
	usable bool
	conn   messageconn.Conn
}

func NewConn(conn messageconn.Conn) *Conn {
	return &Conn{
		usable: true,
		conn:   conn,
	}
}

func (c *Conn) ReadMessage() (msg []byte, err error) {
	return c.conn.ReadMessage()
}

func (c *Conn) WriteMessage(msg []byte) error {
	return c.conn.WriteMessage(msg)
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) Usable() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.usable
}

func (c *Conn) SetUsable() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.usable = true
}

func (c *Conn) SetUnusable() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.usable = false
}

func (c *Conn) ReplaceConn(conn messageconn.Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.conn.Close()
	c.conn = conn
}

// ConnList is a circle linked list of connections.
type ConnList struct {
	lock sync.Mutex

	// list of connections.
	list *ring.Ring
}

// NewList return a new instance of ConnList.
func NewList() *ConnList {
	return &ConnList{
		list: ring.New(0),
	}
}

func (cl *ConnList) Size() int {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.list.Len()
}

func (cl *ConnList) Add(conn *Conn) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	r := &ring.Ring{Value: conn}
	cl.list = r.Link(cl.list)
}

func (cl *ConnList) Remove(c *Conn) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	if cl.list.Len() == 0 {
		return
	}
	if cl.list.Len() == 1 {
		if cl.list.Value.(*Conn) == c {
			cl.list = nil
		}
		return
	}
	for i := cl.list.Len(); i > 0; i-- {
		if cl.list.Value.(*Conn) == c {
			cl.list = cl.list.Prev()
			cl.list.Unlink(1)
			return
		}
		cl.list = cl.list.Next()
	}
}

// Next shift connection to the next one.
func (cl *ConnList) Next() *Conn {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	if cl.list.Len() == 0 {
		return nil
	}
	cl.list = cl.list.Next()
	return cl.list.Value.(*Conn)
}

func (cl *ConnList) ForEach(f func(*Conn)) {
	cl.list.Do(func(v interface{}) {
		f(v.(*Conn))
	})
}

func (cl *ConnList) GetCurConn() *Conn {
	if cl.list.Len() == 0 {
		return nil
	}
	return cl.list.Value.(*Conn)
}
