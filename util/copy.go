package util

import (
	"io"
	"sync"
)

var copyIOBufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 65536)
	},
}

func CopyIO(dst io.WriteCloser, src io.ReadCloser) {
	buf := copyIOBufPool.Get().([]byte)
	defer copyIOBufPool.Put(buf)

	for {
		n, err := src.Read(buf)
		if n > 0 {
			_, err := dst.Write(buf[:n])
			if err != nil {
				src.Close()
				dst.Close()
				return
			}
		}
		if err != nil {
			if err != io.EOF {
				src.Close()
			}
			dst.Close()
			return
		}
	}
}

func CopyBuf(buf []byte) []byte {
	return append(buf[:0:0], buf...)
}
