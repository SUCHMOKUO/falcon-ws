package util

import (
	"testing"
)

func BenchmarkCopyBuf(b *testing.B) {
	buf := make([]byte, 65536)
	for i := 0; i < 65536; i++ {
		buf[i] = byte(RandomUint32())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CopyBuf(buf)
	}
}

func BenchmarkCopy(b *testing.B) {
	buf := make([]byte, 65536)
	for i := 0; i < 65536; i++ {
		buf[i] = byte(RandomUint32())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := make([]byte, 65536)
		copy(c, buf)
	}
}

func BenchmarkCopyFor(b *testing.B) {
	buf := make([]byte, 65536)
	for i := 0; i < 65536; i++ {
		buf[i] = byte(RandomUint32())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := make([]byte, 65536)
		for ii := 0; ii < 65536; ii++ {
			c[ii] = buf[ii]
		}
	}
}

func TestCopyBuf(t *testing.T) {
	buf := []byte{0, 1, 2, 3, 4}
	buf1 := CopyBuf(buf)
	if len(buf1) != len(buf) {
		t.Error("length error")
	}
	for i, n := range buf {
		if buf1[i] != n {
			t.Errorf("buf1[%d] should be %d\n", i, n)
		}
	}
	buf[1] = 100
	if buf1[1] != 1 {
		t.Error("copy error")
	}
}

const testAllocSize = 65536

func allocHeap() {
	buf := make([]byte, testAllocSize)
	buf[testAllocSize-1] = 1
}

func allocStack() {
	var bufArr [testAllocSize]byte
	buf := bufArr[:]
	buf[testAllocSize-1] = 1
}

func BenchmarkAllocateHeap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		allocHeap()
	}
}

func BenchmarkAllocateStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		allocStack()
	}
}
