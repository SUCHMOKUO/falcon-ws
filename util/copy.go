package util

import (
	"io"

	"github.com/SUCHMOKUO/falcon-ws/configs"
)

func CopyIO(dst io.WriteCloser, src io.ReadCloser) {
	// TODO: performance improve. buf will be moved to heap... but why?
	var buf [configs.MaxPackageSize]byte

	for {
		n, err := src.Read(buf[:])
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
