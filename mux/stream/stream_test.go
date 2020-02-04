package stream

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	var id uint32 = 123
	s := New(id)
	done := make(chan bool)
	go func() {
		res, err := ioutil.ReadAll(s)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(res, []byte{0,1,2,3,4}) {
			t.Error(res)
		}
		done <- true
	}()

	s.PutFrame(&Frame{
		Ctl:      DATA,
		StreamId: id,
		Seq:      4,
		Data:     []byte{4},
	})
	s.PutFrame(&Frame{
		Ctl:      FIN,
		StreamId: id,
		Seq:      5,
		Data:     nil,
	})
	s.PutFrame(&Frame{
		Ctl:      DATA,
		StreamId: id,
		Seq:      0,
		Data:     []byte{0},
	})
	s.PutFrame(&Frame{
		Ctl:      DATA,
		StreamId: id,
		Seq:      1,
		Data:     []byte{1},
	})
	s.PutFrame(&Frame{
		Ctl:      DATA,
		StreamId: id,
		Seq:      3,
		Data:     []byte{3},
	})
	s.PutFrame(&Frame{
		Ctl:      DATA,
		StreamId: id,
		Seq:      2,
		Data:     []byte{2},
	})

	<-done
}

func TestWrite(t *testing.T) {
	w := New(1)
	r := New(1)
	done := make(chan bool)

	go func() {
		for {
			f, err := w.GetFrame()
			if err != nil {
				t.Error(err)
			}
			r.PutFrame(f)
			if f.Ctl == FIN {
				break
			}
		}
		done <- true
	}()

	var i byte = 0
	for ; i < 10; i++ {
		w.Write([]byte{i})
	}
	w.Close()

	buf := make([]byte, 1024)
	var curData byte = 0
	for {
		n, err := r.Read(buf)
		if n > 0 {
			if curData != buf[0] {
				t.Error("current data should be:", curData)
			}
			curData++
		}
		if err != nil {
			break
		}
	}

	<-done
}

func BenchmarkFrameSortMap(b *testing.B) {
	buf := make(map[uint16]*Frame, chanBufSize)
	frames := make([]*Frame, 10000)
	for i, _ := range frames {
		frames[i] = &Frame{ Seq: uint16(i) }
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(frames), func(i, j int) {
		frames[i], frames[j] = frames[j], frames[i]
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var nextReadSeq uint16 = 0
		for _, f := range frames {
			buf[f.Seq] = f
			for {
				_, exist := buf[nextReadSeq]
				if !exist {
					break
				}
				delete(buf, nextReadSeq)
				nextReadSeq++
			}
		}
	}
}

func BenchmarkFrameSortArray(b *testing.B) {
	var buf [1<<16]*Frame
	frames := make([]*Frame, 10000)
	for i, _ := range frames {
		frames[i] = &Frame{ Seq: uint16(i) }
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(frames), func(i, j int) {
		frames[i], frames[j] = frames[j], frames[i]
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var nextReadSeq uint16 = 0
		for _, f := range frames {
			buf[f.Seq] = f
			for {
				nextFrame := buf[nextReadSeq]
				if nextFrame == nil {
					break
				}
				buf[nextReadSeq] = nil
				nextReadSeq++
			}
		}
	}
}

func BenchmarkResetBufLoop(b *testing.B) {
	var buf readFrameBuffer

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := range buf {
			buf[j] = nil
		}
	}
}

func BenchmarkResetBufCopy(b *testing.B) {
	var template readFrameBuffer
	var buf readFrameBuffer

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		copy(buf[:], template[:])
	}
}

func BenchmarkResetBufRepeatCopy(b *testing.B) {
	var buf readFrameBuffer

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf[0] = nil
		for bp := 1; bp < 1<<16; bp *= 2 {
			copy(buf[bp:], buf[:bp])
		}
	}
}
