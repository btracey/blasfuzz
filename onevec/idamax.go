// idamaxfuzz provides a fuzzer for the idamax function comparing cgo and native.
package onevec

import (
	"fmt"

	"github.com/btracey/blasfuzz"
	"github.com/gonum/blas/cgo"
	"github.com/gonum/blas/native"
	"github.com/gonum/floats"
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

	alpha, ok := blasfuzz.F64(data)
	if !ok {
		return -1
	}
	data = data[8:]

	str := fmt.Sprintf("Case. N = %v, IncX = %v, x = %#v, alpha = %v", n, incX, x, alpha)

	testIdamax(str, n, x, incX)
	testDnrm2(str, n, x, incX)
	testDasum(str, n, x, incX)
	testDscal(str, n, x, incX, alpha)
	return 0
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
