package util

import (
	"crypto/rand"
	"math"
	"math/big"
	"strconv"
)

func RandomUint32() uint32 {
	res, _ := rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
	return uint32(res.Uint64())
}

func RandomUint32String() string {
	return strconv.FormatUint(uint64(RandomUint32()), 10)
}