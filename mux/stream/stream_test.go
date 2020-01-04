package stream

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestStream_Read(t *testing.T) {
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

func TestStream_Write(t *testing.T) {
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
