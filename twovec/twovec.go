package twovec

import (
	"fmt"

	"github.com/btracey/blasfuzz"
	"github.com/gonum/blas/native"
)

func Fuzz(data []byte) int {
	n, ok := blasfuzz.Int(data, 1)
	if !ok {
		return -1
	}
	data = data[1:]
	incX, ok := blasfuzz.Int(data, 1)
	if !ok {
		return -1
	}
	data = data[1:]
	lenX, ok := blasfuzz.Int(data, 2)
	if !ok {
		return -1
	}
	data = data[2:]
	x, ok := blasfuzz.F64S(data, lenX)
	if !ok {
		return -1
	}
	data = data[lenX*8:]

	incY, ok := blasfuzz.Int(data, 1)
	if !ok {
		return -1
	}
	data = data[1:]
	lenY, ok := blasfuzz.Int(data, 2)
	if !ok {
		return -1
	}
	data = data[2:]
	y, ok := blasfuzz.F64S(data, lenY)
	if !ok {
		return -1
	}
	data = data[lenY*8:]

	alpha, ok := blasfuzz.F64(data)
	if !ok {
		return -1
	}
	data = data[8:]

	str := fmt.Sprintf("Case N = %v\n IncX: %v, x: %v\nIncY: %v, y: %v\n alpha: %v", n, incX, x, incY, y, alpha)
	testDaxpy(str, n, alpha, x, incX, y, incY)
	return 0
}

func testDaxpy(str string, n int, alpha float64, x []float64, incX int, y []float64, incY int) {
	natX := blasfuzz.CloneF64S(x)
	natY := blasfuzz.CloneF64S(y)
	cX := blasfuzz.CloneF64S(x)
	cY := blasfuzz.CloneF64S(y)

	nat := func() { native.Implementation{}.Daxpy(n, alpha, natX, incX, natY, incY) }
	errNative := blasfuzz.CatchPanic(nat)
	c := func() { native.Implementation{}.Daxpy(n, alpha, cX, incX, cY, incY) }
	errC := blasfuzz.CatchPanic(c)

	blasfuzz.SamePanic(str, errC, errNative)
	blasfuzz.SameF64S(str, cX, natX)
	blasfuzz.SameF64S(str, cY, natY)
}
