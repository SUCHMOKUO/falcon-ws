package util

import (
	"io"
	"sync"
)

// buffer pool.
var bufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 1500)
	},
}

func Copy(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()

	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)

	for {
		n, err := src.Read(buf)
		if n > 0 {
			_, err := dst.Write(buf[:n])
			if err != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}