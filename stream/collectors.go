package stream

import (
	"fmt"
	"strings"

	"github.com/BennettJames/ef"
)

// ToMap takes each value in a pair-stream, and turns it into a map where the
// keys are the first value in the pairs, and the values the second.
//
// Note that this cannot handle key collisions - if two pairs have the same `T`
// value, this will panic. Use `StreamToMapMerge` to resolve collisions.
func ToMap[T comparable, U any](srcSt ef.Stream[ef.Pair[T, U]]) map[T]U {
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

// ToMapMerge gathers a pair stream into a map, and resolves any duplicate keys
// using the merge function to combine values.
func ToMapMerge[T comparable, U any](
	srcSt ef.Stream[ef.Pair[T, U]],
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

// Reduce combines all the values in the stream down to one of type `U`. A value
// of `U` type is initialized, and merge is called repeatedly on each element in
// the stream with the current value of `U`. Once the stream is finished, the
// final value is returned.
//
// Example:
//
//	st := StreamOfVals(1, 2, 3)
//	sum := Reduce(st, func(v1, v2 int) int { return v1 + v2 })
func Reduce[T any, U any](
	srcSt ef.Stream[T],
	reduceOp func(total U, val T) U,
) U {
	// todo [bs]: workshop the initialization here
	u := new(U)
	return ReduceInit(srcSt, *u, reduceOp)
}

// ReduceInit combines all the values in the stream down to one of type `U`.
// `merge“ is called repeatedly on each element in the stream with the current
// value of `U`, starting with the provided `initVal`. Once the stream is
// finished, the final value is returned.
//
// Example:
//
//	st := StreamOfVals(1, 2, 3)
//	product := ReduceInit(st, 1, func(v1, v2 int) int { return v1 * v2 })
func ReduceInit[T any, U any](
	srcSt ef.Stream[T],
	initVal U,
	reduceOp func(total U, val T) U,
) U {
	srcSt.Each(func(v T) {
		initVal = reduceOp(initVal, v)
	})
	return initVal
}

// Find searches the stream for a value that matches the provided find operator.
// If a value is found, then it is returned; otherwise an empty optional is.
func Find[T any](
	srcSt ef.Stream[T],
	findOp func(T) bool,
) ef.Opt[T] {
	var foundVal ef.Opt[T]
	srcSt.ExitableEach(func(val T) (advance bool) {
		if findOp(val) {
			foundVal = ef.NewOptValue(val)
			return false
		}
		return true
	})
	return foundVal
}

// Match will return true if any element in the source stream passes the given
// match operator, and false otherwise.
func Match[T any](
	srcSt ef.Stream[T],
	matchOp func(T) bool,
) bool {
	anyMatch := false
	srcSt.ExitableEach(func(val T) (advance bool) {
		if matchOp(val) {
			anyMatch = true
			return false
		}
		return true
	})
	return anyMatch
}

// AnyMatch will return true if every element in the source stream passes the
// given match operator, and false otherwise.
func AllMatch[T any](
	srcSt ef.Stream[T],
	matchOp func(T) bool,
) bool {
	allMatch := true
	srcSt.ExitableEach(func(val T) (advance bool) {
		if !matchOp(val) {
			allMatch = false
			return false
		}
		return true
	})
	return allMatch
}

// Stats calculates the SummaryStats object for a stream of numbers.
func Stats[N ef.Number](srcSt ef.Stream[N]) ef.SummaryStats[N] {
	// note [bs]: possible it'd just be better to make these optionals rather than
	// have default vals for them.
	stats := ef.SummaryStats[N]{
		Min: ef.MaxNumber[N](),
		Max: ef.MinNumber[N](),
	}
	srcSt.Each(func(v N) {
		stats.Size++
		stats.Total += v
		stats.Min = ef.Min(stats.Min, v)
		stats.Max = ef.Max(stats.Max, v)
	})
	if stats.Size > 0 {
		stats.Average = float64(stats.Total) / float64(stats.Size)
	}
	return stats
}

// JoinString combines a stream of strings to a single string, adding `sep`
// between each string.
func JoinString(srcSt ef.Stream[string], sep string) string {
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
