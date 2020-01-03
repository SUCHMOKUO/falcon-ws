package stream

import (
	"bytes"
	"testing"
)

func TestFrame_SerializeAndDeserialize(t *testing.T) {
	ff := []Frame{
		{
			Ctl:      DATA,
			StreamId: 123123,
			Seq:      45,
		},
		{
			Ctl:      FIN,
			StreamId: 842561,
			Seq:      0,
		},
		{
			Ctl:      FIN,
			StreamId: 3141,
			Seq:      0,
			Data:     []byte{1, 2, 3, 4},
		},
	}
	for _, f := range ff {
		s := f.Serialize()
		f2 := DeserializeFrame(s)
		if f2.Ctl != f.Ctl ||
			f2.StreamId != f.StreamId ||
			f2.Seq != f.Seq ||
			!bytes.Equal(f2.Data, f.Data) {
			t.Error(f, f2)
		}
	}
}
