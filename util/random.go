package util

import (
	"crypto/rand"
	"math"
	"math/big"
	"strconv"
)

func RandomUint64() uint64 {
	res, _ := rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
	return res.Uint64()
}

func RandomUint32() uint32 {
	return uint32(RandomUint64())
}

func RandomUintString() string {
	return strconv.FormatUint(RandomUint64(), 10)
}