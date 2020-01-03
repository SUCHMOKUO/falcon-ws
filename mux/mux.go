package mux

import (
	"errors"
	"log"
	"sync"

	"github.com/SUCHMOKUO/falcon-ws/mux/conngroup"
	"github.com/SUCHMOKUO/falcon-ws/mux/stream"
	"github.com/SUCHMOKUO/falcon-ws/util"
)

const (
	chanBufSize = 64
)

var (
	errMuxClosed = errors.New("mux closed")
)

type Mux struct {
	cg            *conngroup.ConnGroup
	streams       sync.Map
	newStreamChan chan *stream.Stream
	closed        bool
}

func New(cg *conngroup.ConnGroup) *Mux {
	m := new(Mux)
	m.cg = cg
	m.newStreamChan = make(chan *stream.Stream, chanBufSize)
	go m.streamsDataReceiver()
	return m
}

func (m *Mux) Close() {
	m.closed = true
	m.cg.Close()
	close(m.newStreamChan)
}

func (m *Mux) NewStream() (*stream.Stream, error) {
	if m.closed {
		return nil, errMuxClosed
	}
	return m.newStream(util.RandomUint32()), nil
}

func (m *Mux) newStream(id uint32) *stream.Stream {
	s := stream.New(id)
	m.streams.Store(s.Id, s)
	s.OnFullClose = m.onStreamClose
	go m.streamDataSender(s)
	return s
}

func (m *Mux) Accept() (*stream.Stream, error) {
	s, ok := <-m.newStreamChan
	if ok {
		return s, nil
	}
	return nil, errMuxClosed
}

func (m *Mux) streamDataSender(s *stream.Stream) {
	for {
		f, err := s.GetFrame()
		if err != nil {
			return
		}
		m.cg.WriteMessage(f.Serialize())
		log.Printf("[Stream] %d, type: %d, seq: %d, send %d bytes.\n", f.StreamId, f.Ctl, f.Seq, len(f.Data))
	}
}

func (m *Mux) streamsDataReceiver() {
	for {
		msg, err := m.cg.ReadMessage()
		if err != nil {
			return
		}
		f := stream.DeserializeFrame(msg)
		if f == nil {
			continue
		}
		v, ok := m.streams.Load(f.StreamId)
		var s *stream.Stream
		if ok {
			s = v.(*stream.Stream)
		} else {
			s = m.newStream(f.StreamId)
			m.newStreamChan <- s
		}
		s.PutFrame(f)
		log.Printf("[Stream] %d, type: %d, seq: %d, received %d bytes.\n", s.Id, f.Ctl, f.Seq, len(f.Data))
	}
}

func (m *Mux) onStreamClose(s *stream.Stream) {
	m.streams.Delete(s.Id)
}
