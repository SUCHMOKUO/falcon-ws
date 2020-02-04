package util

import (
	"io"
)

func CopyIO(dst io.WriteCloser, src io.ReadCloser) {
	var bufArr [10240]byte
	buf := bufArr[:]

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
