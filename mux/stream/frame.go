package stream

import (
	"encoding/binary"
)

// control message.
const (
	DATA = iota
	FIN
)

type Frame struct {
	Ctl      byte
	StreamId uint32
	Seq      uint16
	Data     []byte
}

func (f *Frame) Serialize() []byte {
	buf := make([]byte, len(f.Data) + 7)
	buf[0] = f.Ctl
	binary.BigEndian.PutUint32(buf[1:5], f.StreamId)
	binary.BigEndian.PutUint16(buf[5:7], f.Seq)
	copy(buf[7:], f.Data)
	return buf
}

func DeserializeFrame(rawFrame []byte) *Frame {
	msgType := rawFrame[0]
	frame := new(Frame)

	switch msgType {
	case DATA:
		frame.Ctl = DATA
	case FIN:
		frame.Ctl = FIN
	default:
		return nil
	}

	rawFrameLen := len(rawFrame)
	if rawFrameLen >= 5 {
		frame.StreamId = binary.BigEndian.Uint32(rawFrame[1:5])
	}
	if rawFrameLen >= 7 {
		frame.Seq = binary.BigEndian.Uint16(rawFrame[5:7])
	}
	if rawFrameLen > 7 {
		frame.Data = rawFrame[7:]
	}

	return frame
}
