package ef

import (
	"fmt"
	"strings"
)

// StreamToMap takes each value in a pair-stream, and turns it into a map where
// the keys are the first value in the pairs, and the values the second.
//
// Note that this cannot handle key collisions - if two pairs have the same `T`
// value, this will panic. Use `StreamToMapMerge` to resolve collisions.
func StreamToMap[T comparable, U any](srcSt Stream[Pair[T, U]]) map[T]U {
	m := make(map[T]U)
	EachPair(srcSt, func(t T, u U) {
		if existing, exists := m[t]; !exists {
			m[t] = u
		} else {
			// todo [bs]: probably want a custom error for this
			panic(fmt.Errorf(
				"StreamToMap: duplicate values found for key '%v' - ['%v', '%v']",
				t, u, existing))
		}
		m[t] = u
	})
	return m
}

// StreamToMapMerge gathers a pair stream into a map, and resolves any duplicate
// keys using the merge function to combine values.
func StreamToMapMerge[T comparable, U any](
	srcSt Stream[Pair[T, U]],
	mergeOp func(key T, val1, val2 U) U,
) map[T]U {
	m := make(map[T]U)
	EachPair(srcSt, func(key T, value U) {
		if existing, exists := m[key]; !exists {
			m[key] = value
		} else {
			m[key] = mergeOp(key, existing, value)
		}
		m[key] = value
	})
	return m
}

// StreamReduce combines all the values in the stream down to one of type `U`. A
// value of `U` type is initialized, and merge is called repeatedly on each
// element in the stream with the current value of `U`. Once the stream is
// finished, the final value is returned.
//
// Example:
//
//	st := StreamOfVals(1, 2, 3)
//	sum := StreamReduce(st, func(v1, v2 int) int { return v1 + v2 })
func StreamReduce[T any, U any](
	srcSt Stream[T],
	reduceOp func(total U, val T) U,
) U {
	// todo [bs]: workshop the initialization here
	u := new(U)
	return StreamReduceInit(srcSt, *u, reduceOp)
}

// StreamReduceInit combines all the values in the stream down to one of type
// `U`. `merge“ is called repeatedly on each element in the stream with the
// current value of `U`, starting with the provided `initVal`. Once the stream
// is finished, the final value is returned.
//
// Example:
//
//	st := StreamOfVals(1, 2, 3)
//	product := StreamReduceInit(st, 1, func(v1, v2 int) int { return v1 * v2 })
func StreamReduceInit[T any, U any](
	srcSt Stream[T],
	initVal U,
	reduceOp func(total U, val T) U,
) U {
	srcSt.Each(func(v T) {
		initVal = reduceOp(initVal, v)
	})
	return initVal
}

// StreamReduce combines all the values in the pair-stream down to one of type `V`. A
// value of `V` type is initialized, and merge is called repeatedly on each
// element in the stream with the current value of `V`. Once the stream is
// finished, the final value is returned.
//
// Example:
//
//	st := StreamOfMap(map[string]int{
//	  "a": 1,
//	  "b": 2,
//	  "c": 3,
//	})
//	sum := StreamReduce(st, func(total int, k string, v int) int { return total + v })
func PStreamReduce[T any, U any, V any](
	srcSt Stream[Pair[T, U]],
	reduceOp func(total V, first T, second U) V,
) V {
	v := new(V)
	return PStreamReduceInit(srcSt, *v, reduceOp)
}

// PStreamReduceInit combines all the values in the pair-stream down to one of
// type `V`. `merge“ is called repeatedly on each element in the stream with
// the current value of `V`, starting with the provided `initVal`. Once the
// stream is finished, the final value is returned.
//
// Example:
//
//	st := StreamOfMap(map[string]int{
//	  "a": 1,
//	  "b": 2,
//	  "c": 3,
//	})
//	product := PStreamReduceInit(st, 1, func(total int, key string, val int) int {
//	  return total * val
//	})
func PStreamReduceInit[T any, U any, V any](
	srcSt Stream[Pair[T, U]],
	initVal V,
	reduceOp func(total V, first T, second U) V,
) V {
	EachPair(srcSt, func(v1 T, v2 U) {
		initVal = reduceOp(initVal, v1, v2)
	})
	return initVal
}

// StreamFind searches the stream for a value that matches the provided find
// operator. If a value is found, then it is returned; otherwise an empty
// optional is.
func StreamFind[T any](
	srcSt Stream[T],
	findOp func(T) bool,
) Opt[T] {
	var foundVal Opt[T]
	srcSt.srcIter.iterate(func(val T) (advance bool) {
		if findOp(val) {
			foundVal = NewOptValue(val)
			return false
		}
		return true
	})
	return foundVal
}

// PStreamFind searches the pair-stream for a value that matches the provided
// find operator. If a value is found, then it is returned; otherwise an empty
// optional is.
func PStreamFind[T, U any](
	srcSt Stream[Pair[T, U]],
	findOp func(T, U) bool,
) Opt[Pair[T, U]] {
	return StreamFind(srcSt, func(p Pair[T, U]) bool {
		return findOp(p.Get())
	})
}

// StreamAnyMatch will return true if any element in the source stream passes
// the given match operator, and false otherwise.
func StreamAnyMatch[T any](
	srcSt Stream[T],
	matchOp func(T) bool,
) bool {
	anyMatch := false
	srcSt.srcIter.iterate(func(val T) (advance bool) {
		if matchOp(val) {
			anyMatch = true
			return false
		}
		return true
	})
	return anyMatch
}

// PStreamAnyMatch will return true if any element in the source pair stream
// passes the given match operator, and false otherwise.
func PStreamAnyMatch[T, U any](
	srcSt Stream[Pair[T, U]],
	matchOp func(T, U) bool,
) bool {
	return StreamAnyMatch(srcSt, func(p Pair[T, U]) bool {
		return matchOp(p.Get())
	})
}

// StreamAnyMatch will return true if every element in the source stream passes
// the given match operator, and false otherwise.
func StreamAllMatch[T any](
	srcSt Stream[T],
	matchOp func(T) bool,
) bool {
	allMatch := true
	srcSt.srcIter.iterate(func(val T) (advance bool) {
		if !matchOp(val) {
			allMatch = false
			return false
		}
		return true
	})
	return allMatch
}

// PStreamAllMatch will return true if every element in the source pair stream
// passes the given match operator, and false otherwise.
func PStreamAllMatch[T, U any](
	srcSt Stream[Pair[T, U]],
	matchOp func(T, U) bool,
) bool {
	return StreamAllMatch(srcSt, func(p Pair[T, U]) bool {
		return matchOp(p.Get())
	})
}

// SummaryStats contains a set of data about the values in a stream of numbers.
//
// Note that this is not safe with overflow - if the sum exceeds the number
// type, then overflow will occur and total / average will not be accurate.
type SummaryStats[N Number] struct {
	Average  float64
	Size     int
	Total    N
	Min, Max N
}

// StreamStats calculates the SummaryStats object for a stream of numbers.
func StreamStats[N Number](srcSt Stream[N]) SummaryStats[N] {
	// note [bs]: possible it'd just be better to make these optionals rather than
	// have default vals for them.
	stats := SummaryStats[N]{
		Min: MaxNumber[N](),
		Max: MinNumber[N](),
	}
	srcSt.Each(func(v N) {
		stats.Size++
		stats.Total += v
		stats.Min = Min(stats.Min, v)
		stats.Max = Max(stats.Max, v)
	})
	if stats.Size > 0 {
		stats.Average = float64(stats.Total) / float64(stats.Size)
	}
	return stats
}

// StreamJoinString combines a stream of strings to a single string, adding
// `sep` between each string.
func StreamJoinString(srcSt Stream[string], sep string) string {
	var sb strings.Builder
	first := true
	srcSt.Each(func(v string) {
		if !first {
			sb.WriteString(sep)
		} else {
			first = false
		}
		sb.WriteString(v)
	})
	return sb.String()
}
