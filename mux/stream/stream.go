package stream

import (
	"errors"
	"io"

	"github.com/SUCHMOKUO/falcon-ws/util"
)

const (
	chanBufSize = 64
)

var (
	errWriteToClosedStream = errors.New("write to closed stream")
	errStreamClosed        = errors.New("stream closed")
)

// Stream is a abstract connection beyond the connection group.
// Stream is NOT concurrency-safe.
type Stream struct {
	Id          uint32
	OnFullClose func(stream *Stream)

	nextReadSeq  uint16
	nextWriteSeq uint16

	// used for caching read frames that out of order.
	readFrameBuf map[uint16]*Frame

	// used for reading unordered frames.
	receiveChan chan *Frame

	// used for reading ordered frame.
	readChan chan *Frame

	// used for sending frames.
	sendChan chan *Frame

	readClosed  bool
	writeClosed bool
}

func New(id uint32) *Stream {
	s := new(Stream)
	s.Id = id
	s.readFrameBuf = make(map[uint16]*Frame, chanBufSize)
	s.readChan = make(chan *Frame, chanBufSize)
	s.receiveChan = make(chan *Frame, chanBufSize)
	s.sendChan = make(chan *Frame, chanBufSize)
	go s.receivedFrameSorter()
	return s
}

func (s *Stream) WriteClosed() bool {
	return s.writeClosed
}

func (s *Stream) ReadClosed() bool {
	return s.readClosed
}

func (s *Stream) FullClosed() bool {
	return s.writeClosed && s.readClosed
}

func (s *Stream) Read(p []byte) (n int, err error) {
	if s.readClosed {
		return 0, io.EOF
	}
	f := <-s.readChan
	if f.Ctl == FIN {
		s.closeRead()
		err = io.EOF
	}
	n = copy(p, f.Data)
	return n, err
}

func (s *Stream) Write(p []byte) (n int, err error) {
	if s.writeClosed {
		return 0, errWriteToClosedStream
	}
	pLen := len(p)
	data := util.CopyBuf(p)
	f := &Frame{
		Ctl:      DATA,
		StreamId: s.Id,
		Seq:      s.nextWriteSeq,
		Data:     data,
	}
	s.nextWriteSeq++
	s.sendChan <- f
	return pLen, nil
}

func (s *Stream) Close() error {
	if s.writeClosed {
		return errStreamClosed
	}
	fin := &Frame{
		Ctl:      FIN,
		StreamId: s.Id,
		Seq:      s.nextWriteSeq,
	}
	s.sendChan <- fin
	s.closeWrite()
	return nil
}

func (s *Stream) PutFrame(f *Frame) {
	if s.readClosed {
		return
	}
	s.receiveChan <- f
}

func (s *Stream) GetFrame() (*Frame, error) {
	f, ok := <-s.sendChan
	if !ok {
		return nil, errStreamClosed
	}
	return f, nil
}

func (s *Stream) closeRead() {
	if s.readClosed {
		return
	}
	s.readClosed = true
	if s.FullClosed() {
		s.clear()
	}
}

func (s *Stream) closeWrite() {
	if s.writeClosed {
		return
	}
	s.writeClosed = true
	if s.FullClosed() {
		s.clear()
	}
}

func (s *Stream) clear() {
	close(s.sendChan)
	close(s.readChan)
	close(s.receiveChan)
	if s.OnFullClose != nil {
		s.OnFullClose(s)
	}
}

func (s *Stream) receivedFrameSorter() {
	for {
		f, ok := <-s.receiveChan
		if !ok {
			return
		}
		s.readFrameBuf[f.Seq] = f
		for {
			nextFrame, exist := s.readFrameBuf[s.nextReadSeq]
			if !exist {
				break
			}
			delete(s.readFrameBuf, s.nextReadSeq)
			s.nextReadSeq++
			s.readChan <- nextFrame
		}
	}
}
