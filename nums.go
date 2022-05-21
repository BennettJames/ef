package ef

import "math"

// SingedInteger is an interface union of all signed integer primitives.
type SingedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedInteger is an interface union of all unsigned integer primitives.
type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is an interface union of all integer primitives.
type Integer interface {
	SingedInteger | UnsignedInteger
}

// Float is an interfaces union of all float primitives.
type Float interface {
	float32 | float64
}

// Complex is an interface union of all explicitly complex primitives.
type Complex interface {
	complex64 | complex128
}

// Number is an interface union of all real numbers.
//
// Authors note: this is called "Number" rather than "RealNumber" mostly as it
// is shorter, and order-able numbers are used much more often than complex.
type Number interface {
	Integer | Float
}

// AllNumber is an inferace union of all complex and real number types.
type AllNumber interface {
	Number | Complex
}

// Order returns the two given numbers so the first is the lower value.
func Order[N Number](v1, v2 N) (low N, high N) {
	if v1 < v2 {
		return v1, v2
	} else {
		return v2, v1
	}
}

// Range returns a stream that consists of values from the start until the end.
// The end is exclusive - that is, the stream consists of integers less than the
// end.
func Range[I Integer](start, end I) Stream[I] {
	return Stream[I]{
		src: &iterInts[I]{
			start: start,
			end:   end,
		},
	}
}

// iterInts is a simple iterator that supports the range stream.
type iterInts[I Integer] struct {
	start, end I
	index      I
}

func (i *iterInts[I]) Next() Opt[I] {
	v := i.start + i.index
	if v > i.end {
		return Opt[I]{}
	}
	i.index += 1
	return OptOf(v)
}

// Min returns the lower of the two values.
func Min[N Number](v1, v2 N) N {
	if v1 <= v2 {
		return v1
	}
	return v2
}

// Max returns the higher of the two values.
func Max[N Number](v1, v2 N) N {
	if v1 >= v2 {
		return v1
	}
	return v2
}

// MinNumber returns the minimum possible number for the given number type.
func MinNumber[N Number]() N {
	switch any(*new(N)).(type) {
	case int:
		x := int(math.MinInt)
		return N(x)
	case int8:
		x := math.MinInt8
		return N(x)
	case int16:
		x := (math.MinInt16)
		return N(x)
	case int32:
		x := math.MinInt32
		return N(x)
	case int64:
		x := math.MinInt64
		return N(x)
	case uint:
		return N(0)
	case uint8:
		return N(0)
	case uint16:
		return N(0)
	case uint32:
		return N(0)
	case uint64:
		return N(0)
	case uintptr:
		return N(0)
	case float32:
		x := float32(math.Inf(-1))
		return N(x)
	case float64:
		x := float64(math.Inf(-1))
		return N(x)
	default:
		panic("unreachable")
	}
}

// MaxNumber returns the maximum possible value for the given number type.
func MaxNumber[N Number]() N {
	switch any(*new(N)).(type) {
	case int:
		x := int(math.MaxInt)
		return N(x)
	case int8:
		x := int8(math.MaxInt8)
		return N(x)
	case int16:
		x := int16(math.MaxInt16)
		return N(x)
	case int32:
		x := int32(math.MaxInt32)
		return N(x)
	case int64:
		x := int64(math.MaxInt64)
		return N(x)
	case uint:
		x := uint(math.MaxUint)
		return N(x)
	case uint8:
		x := uint8(math.MaxUint8)
		return N(x)
	case uint16:
		x := uint16(math.MaxUint16)
		return N(x)
	case uint32:
		x := uint32(math.MaxUint32)
		return N(x)
	case uint64:
		x := uint64(math.MaxUint64)
		return N(x)
	case uintptr:
		x := uintptr(math.MaxUint64)
		return N(x)
	case float32:
		x := float32(math.Inf(1))
		return N(x)
	case float64:
		x := float64(math.Inf(1))
		return N(x)
	default:
		panic("unreachable")
	}
}
