package twovec

import (
	"fmt"

	"github.com/btracey/blasfuzz"
	"github.com/gonum/blas/cgo"
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

	beta, ok := blasfuzz.F64(data)
	if !ok {
		return -1
	}
	data = data[8:]

	str := fmt.Sprintf("Case N = %v\n IncX: %v, x: %v\nIncY: %v, y: %v\n alpha: %v", n, incX, x, incY, y, alpha)
	testDaxpy(str, n, alpha, x, incX, y, incY)
	testDcopy(str, n, x, incX, y, incY)
	testDdot(str, n, x, incX, y, incY)
	testDswap(str, n, x, incX, y, incY)
	testDrot(str, n, x, incX, y, incY, alpha, beta)
	return 0
}

func testDaxpy(str string, n int, alpha float64, x []float64, incX int, y []float64, incY int) {
	natX := blasfuzz.CloneF64S(x)
	natY := blasfuzz.CloneF64S(y)
	cX := blasfuzz.CloneF64S(x)
	cY := blasfuzz.CloneF64S(y)

	nat := func() { native.Implementation{}.Daxpy(n, alpha, natX, incX, natY, incY) }
	errNative := blasfuzz.CatchPanic(nat)
	c := func() { cgo.Implementation{}.Daxpy(n, alpha, cX, incX, cY, incY) }
	errC := blasfuzz.CatchPanic(c)
	if n == 0 {
		str2 := fmt.Sprintf("cerr = %v, naterr = %v", errC, errNative)
		panic(str2)
	}
	blasfuzz.SamePanic(str, errC, errNative)
	blasfuzz.SameF64S(str, cX, natX)
	blasfuzz.SameF64S(str, cY, natY)
}

func testDcopy(str string, n int, x []float64, incX int, y []float64, incY int) {
	natX := blasfuzz.CloneF64S(x)
	natY := blasfuzz.CloneF64S(y)
	cX := blasfuzz.CloneF64S(x)
	cY := blasfuzz.CloneF64S(y)

	nat := func() { native.Implementation{}.Dcopy(n, natX, incX, natY, incY) }
	errNative := blasfuzz.CatchPanic(nat)
	c := func() { cgo.Implementation{}.Dcopy(n, cX, incX, cY, incY) }
	errC := blasfuzz.CatchPanic(c)

	blasfuzz.SamePanic(str, errC, errNative)
	blasfuzz.SameF64S(str, cX, natX)
	blasfuzz.SameF64S(str, cY, natY)
}

func testDswap(str string, n int, x []float64, incX int, y []float64, incY int) {
	natX := blasfuzz.CloneF64S(x)
	natY := blasfuzz.CloneF64S(y)
	cX := blasfuzz.CloneF64S(x)
	cY := blasfuzz.CloneF64S(y)

	nat := func() { native.Implementation{}.Dswap(n, natX, incX, natY, incY) }
	errNative := blasfuzz.CatchPanic(nat)
	c := func() { cgo.Implementation{}.Dswap(n, cX, incX, cY, incY) }
	errC := blasfuzz.CatchPanic(c)

	blasfuzz.SamePanic(str, errC, errNative)
	blasfuzz.SameF64S(str, cX, natX)
	blasfuzz.SameF64S(str, cY, natY)
}

func testDdot(str string, n int, x []float64, incX int, y []float64, incY int) {
	natX := blasfuzz.CloneF64S(x)
	natY := blasfuzz.CloneF64S(y)
	cX := blasfuzz.CloneF64S(x)
	cY := blasfuzz.CloneF64S(y)

	var natAns float64
	nat := func() { natAns = native.Implementation{}.Ddot(n, natX, incX, natY, incY) }
	errNative := blasfuzz.CatchPanic(nat)
	var cAns float64
	c := func() { cAns = cgo.Implementation{}.Ddot(n, cX, incX, cY, incY) }
	errC := blasfuzz.CatchPanic(c)

	blasfuzz.SamePanic(str, errC, errNative)
	blasfuzz.SameF64S(str, cX, natX)
	blasfuzz.SameF64S(str, cY, natY)
	blasfuzz.SameF64Approx(str, cAns, natAns, 1e-13, 1e-13)
}

func testDrot(str string, n int, x []float64, incX int, y []float64, incY int, c, s float64) {
	natX := blasfuzz.CloneF64S(x)
	natY := blasfuzz.CloneF64S(y)
	cX := blasfuzz.CloneF64S(x)
	cY := blasfuzz.CloneF64S(y)

	nat := func() { native.Implementation{}.Drot(n, natX, incX, natY, incY, c, s) }
	errNative := blasfuzz.CatchPanic(nat)
	cFunc := func() { cgo.Implementation{}.Drot(n, cX, incX, cY, incY, c, s) }
	errC := blasfuzz.CatchPanic(cFunc)

	blasfuzz.SamePanic(str, errC, errNative)
	blasfuzz.SameF64S(str, cX, natX)
	blasfuzz.SameF64S(str, cY, natY)
}
