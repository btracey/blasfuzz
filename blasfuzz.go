package blasfuzz

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/gonum/floats"
)

// GetInt gets an integer with the given number of bytes
func Int(data []byte, b int) (n int, ok bool) {
	if len(data) < b {
		return 0, false
	}
	if b == 1 {
		return int(data[0]), true
	}
	if b == 2 {
		return int(binary.LittleEndian.Uint16(data[:2:2])), true
	}
	panic("not coded")
}

func F64(data []byte) (float64, bool) {
	if len(data) < 8 {
		return math.NaN(), false
	}
	uint64 := binary.LittleEndian.Uint64(data[:8:8])
	float64 := math.Float64frombits(uint64)
	return float64, true
}

func F64S(data []byte, l int) ([]float64, bool) {
	var ok bool
	x := make([]float64, l)
	for i := range x {
		x[i], ok = F64(data)
		if !ok {
			return nil, false
		}
		data = data[8:]
	}
	return x, true
}

// Panics returns the error if panics
func CatchPanic(f func()) (err interface{}) {
	defer func() {
		err = recover()
	}()
	f()
	return
}

// SameError checks that the two errors are the same if either of them are non-nil.
func SamePanic(str string, c, native interface{}) {
	if c != nil && native == nil {
		panic(fmt.Sprintf("Case %s: c panics, native doesn't: %v", str, c))
	}
	if c == nil && native != nil {
		panic(fmt.Sprintf("Case %s: native panics, c doesn't: %v", str, native))
	}
	if c != native {
		panic(fmt.Sprintf("Case %s: Error mismatch.\nC is: %v\nNative is: %v", str, c, native))
	}
}

func CloneF64S(x []float64) []float64 {
	n := make([]float64, len(x))
	copy(n, x)
	return n
}

func SameInt(str string, c, native int) {
	if c != native {
		panic(fmt.Sprintf("Case %s: Int mismatch. c = %v, native = %v.", str, c, native))
	}
}

func SameF64Approx(str string, c, native, absTol, relTol float64) {
	if math.IsNaN(c) && math.IsNaN(native) {
		return
	}
	if !floats.EqualWithinAbsOrRel(c, native, absTol, relTol) {
		cb := math.Float64bits(c)
		nb := math.Float64bits(native)
		same := floats.EqualWithinAbsOrRel(c, native, absTol, relTol)
		panic(fmt.Sprintf("Case %s: Float64 mismatch. c = %v, native = %v\n cb: %v, nb: %v\n%v,%v,%v", str, c, native, cb, nb, same, absTol, relTol))
	}
}

func SameF64S(str string, c, native []float64) {
	if !floats.Same(c, native) {
		panic(fmt.Sprintf("Case %s: []float64 mismatch. c = %v, native = %v.", str, c, native))
	}
}
