package ef

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

type Float interface {
	float32 | float64
}

type Complex interface {
	complex64 | complex128
}

type Number interface {
	Integer | Float
}

type AllNumber interface {
	Number | Complex
}

type iterInts[I Integer] struct {
	start, end I
	index      I
}

func order[N Number](v1, v2 N) (N, N) {

	// note [bs]: generally I don't think these sorts of number types really
	// belong in the same package, but for the sake of experimentation eh why not.
	if v1 < v2 {
		return v1, v2
	} else {
		return v2, v1
	}
}

func Range[I Integer](start, end I) Stream[I] {
	return Stream[I]{
		src: &iterInts[I]{
			start: start,
			end:   end,
		},
	}
}

func (i *iterInts[I]) Next() Opt[I] {
	v := i.start + i.index
	if v > i.end {
		return Opt[I]{}
	}
	i.index += 1
	return NewOpt(v)
}

func Min[N Number](v1, v2 N) N {
	if v1 <= v2 {
		return v1
	}
	return v2
}

func Max[N Number](v1, v2 N) N {
	if v1 >= v2 {
		return v1
	}
	return v2
}

const MaxUint = ^uint(0)
const MinUint = 0

// int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr

// float32 | float64

const (
	minInt int = -maxInt - 1
	maxInt int = int(^uint(0) >> 1)

	minInt8 int8 = -maxInt8 - 1
	maxInt8 int8 = int8(^uint8(0) >> 1)

	minInt16 int16 = -maxInt16 - 1
	maxInt16 int16 = int16(^uint16(0) >> 1)

	minInt32 int32 = -maxInt32 - 1
	maxInt32 int32 = int32(^uint32(0) >> 1)

	minInt64 int64 = -maxInt64 - 1
	maxInt64 int64 = int64(^uint64(0) >> 1)

	minUint uint = 0
	maxUint uint = ^uint(0)

	minUint8 uint8 = 0
	maxUint8 uint8 = ^uint8(0)

	minUint16 uint16 = 0
	maxUint16 uint16 = ^uint16(0)

	minUint32 uint32 = 0
	maxUint32 uint32 = ^uint32(0)

	minUint64 uint64 = 0
	maxUint64 uint64 = ^uint64(0)
)
