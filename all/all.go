// all fuzzes all of the BLAS functions
package all

import (
	"fmt"

	"github.com/btracey/blasfuzz"
	"github.com/gonum/blas"
	"github.com/gonum/blas/cgo"
	"github.com/gonum/blas/native"
	"github.com/gonum/floats"
)

func Fuzz(data []byte) int {
	n, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]

	// Construct slice 1
	incX, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]
	lenX, ok := blasfuzz.Int(data, 2)
	if !ok {
		return 0
	}
	data = data[2:]
	x, ok := blasfuzz.F64S(data, lenX)
	if !ok {
		return 0
	}
	data = data[lenX*8:]

	// Construct slice 2
	incY, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]
	lenY, ok := blasfuzz.Int(data, 2)
	if !ok {
		return 0
	}
	data = data[2:]
	y, ok := blasfuzz.F64S(data, lenY)
	if !ok {
		return 0
	}
	data = data[lenY*8:]

	// Construct matrix 1
	m1, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]
	n1, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]
	ld1, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]
	lenA := ld1*m1 + n1
	a, ok := blasfuzz.F64S(data, lenA)
	if !ok {
		return 0
	}
	data = data[lenA*8:]

	// Construct matrix 2
	m2, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]
	n2, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]
	ld2, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]
	lenB := ld2*m2 + n2
	b, ok := blasfuzz.F64S(data, lenB)
	if !ok {
		return 0
	}
	data = data[lenB*8:]

	// Construct matrix 3
	m3, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]
	n3, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]
	ld3, ok := blasfuzz.Int(data, 1)
	if !ok {
		return 0
	}
	data = data[1:]
	lenC := ld3*m3 + n3
	c, ok := blasfuzz.F64S(data, lenC)
	if !ok {
		return 0
	}
	data = data[lenC*8:]

	// Generate a couple of parameters and booleans
	nParams := 8
	params, ok := blasfuzz.F64S(data, nParams)
	if !ok {
		return 0
	}
	data = data[nParams*8:]

	// Generate a couple of booleans
	bools, ok := blasfuzz.Bools(data)
	if !ok {
		return 0
	}
	data = data[1:]

	// Generate an integer
	iParam, ok := blasfuzz.Int(data, 1)
	if !ok {
		return -1
	}
	data = data[1:]

	_, _, _ = a, b, c
	_ = bools

	// Test the functions
	level1Test(n, x, lenX, incX, y, lenY, incY, params, iParam)

	return 1
}

func level1Test(n int, x []float64, lenX, incX int, y []float64, lenY, incY int, params []float64, iParam int) {
	alpha := params[0]
	beta := params[1]

	flag := iParam
	if flag < 0 {
		flag = -flag
	}
	flag = flag%4 - 2

	drotm := blas.DrotmParams{
		blas.Flag(flag),
		[4]float64{params[0], params[1], params[2], params[3]},
	}

	str1 := fmt.Sprintf("Case. N = %v, IncX = %v, x = %#v, alpha = %v", n, incX, x, alpha)
	str2 := fmt.Sprintf("Case N = %v\n IncX: %v, x: %v\nIncY: %v, y: %v\n alpha: %v", n, incX, x, incY, y, alpha)
	str3 := fmt.Sprintf("Case N = %v\n IncX: %v, x: %v\nIncY: %v, y: %v\n alpha: %v beta: %v", n, incX, x, incY, y, alpha, beta)
	str4 := fmt.Sprintf("Case. N = %v, IncX = %v, x = %#v, drotm = %v", n, incX, x, drotm)

	testDrot(str3, n, x, incX, y, incY, alpha, beta)
	testDrotm(str4, n, x, incX, y, incY, drotm)
	testDswap(str2, n, x, incX, y, incY)
	testDscal(str1, n, x, incX, alpha)
	testDcopy(str2, n, x, incX, y, incY)
	testDaxpy(str2, n, alpha, x, incX, y, incY)
	testDdot(str2, n, x, incX, y, incY)
	testDnrm2(str1, n, x, incX)
	testDasum(str1, n, x, incX)
	testIdamax(str1, n, x, incX)
}

func testIdamax(str string, n int, x []float64, incX int) {
	var natAns int
	nat := func() { natAns = native.Implementation{}.Idamax(n, x, incX) }
	errNative := blasfuzz.CatchPanic(nat)

	cx := blasfuzz.CloneF64S(x)
	var cAns int
	c := func() { cAns = cgo.Implementation{}.Idamax(n, cx, incX) }
	errC := blasfuzz.CatchPanic(c)

	blasfuzz.SamePanic(str, errC, errNative)
	blasfuzz.SameF64S(str, cx, x)
	// Known issue: If the slice contains NaN the answer may vary
	if !floats.HasNaN(x) {
		blasfuzz.SameInt(str, cAns, natAns)
	}
}

func testDnrm2(str string, n int, x []float64, incX int) {
	var natAns float64
	nat := func() { natAns = native.Implementation{}.Dnrm2(n, x, incX) }
	errNative := blasfuzz.CatchPanic(nat)

	cx := blasfuzz.CloneF64S(x)
	var cAns float64
	c := func() { cAns = cgo.Implementation{}.Dnrm2(n, cx, incX) }
	errC := blasfuzz.CatchPanic(c)

	blasfuzz.SamePanic(str, errC, errNative)
	blasfuzz.SameF64S(str, cx, x)
	blasfuzz.SameF64Approx(str, cAns, natAns, 1e-13, 1e-13)
}

func testDasum(str string, n int, x []float64, incX int) {
	var natAns float64
	nat := func() { natAns = native.Implementation{}.Dasum(n, x, incX) }
	errNative := blasfuzz.CatchPanic(nat)

	cx := blasfuzz.CloneF64S(x)
	var cAns float64
	c := func() { cAns = cgo.Implementation{}.Dasum(n, cx, incX) }
	errC := blasfuzz.CatchPanic(c)

	blasfuzz.SamePanic(str, errC, errNative)
	blasfuzz.SameF64S(str, cx, x)
	blasfuzz.SameF64Approx(str, cAns, natAns, 1e-13, 1e-13)
}

func testDscal(str string, n int, x []float64, incX int, alpha float64) {
	natX := blasfuzz.CloneF64S(x)
	nat := func() { native.Implementation{}.Dscal(n, alpha, natX, incX) }
	errNative := blasfuzz.CatchPanic(nat)

	cx := blasfuzz.CloneF64S(x)
	c := func() { cgo.Implementation{}.Dscal(n, alpha, cx, incX) }
	errC := blasfuzz.CatchPanic(c)

	blasfuzz.SamePanic(str, errC, errNative)
	blasfuzz.SameF64S(str, cx, natX)
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
	/*
		if n == 0 {
			str2 := fmt.Sprintf("cerr = %v, naterr = %v", errC, errNative)
			panic(str2)
		}
	*/
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

func testDrotm(str string, n int, x []float64, incX int, y []float64, incY int, param blas.DrotmParams) {
	natX := blasfuzz.CloneF64S(x)
	natY := blasfuzz.CloneF64S(y)
	cX := blasfuzz.CloneF64S(x)
	cY := blasfuzz.CloneF64S(y)
	nat := func() { native.Implementation{}.Drotm(n, natX, incX, natY, incY, param) }
	errNative := blasfuzz.CatchPanic(nat)
	cFunc := func() { cgo.Implementation{}.Drotm(n, cX, incX, cY, incY, param) }
	errC := blasfuzz.CatchPanic(cFunc)

	blasfuzz.SamePanic(str, errC, errNative)
	blasfuzz.SameF64S(str, cX, natX)
	blasfuzz.SameF64S(str, cY, natY)
}
