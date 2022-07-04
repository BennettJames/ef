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

// StreamMap transforms each value in the input stream into a new value with
// the provided function, and returns a new stream with the result.
func StreamMap[T, U any](srcSt Stream[T], mapOp func(v T) U) Stream[U] {

	// so - this feels like it ought be simplifiable, or at least captured
	// with a more consistent pattern. Rough idea: most fn iter's have at least
	// some basic consistent patterns. In practice, they often end up doing
	// something very similar to the original

	return Stream[U]{
		srcIter: &fnIter[U]{
			fn: func(wrappedFn func(U) bool) {
				srcSt.srcIter.iterate(func(val T) bool {
					return wrappedFn(mapOp(val))
				})
			},
		},
	}
}

func StreamMap2[T, U any, S Streamable[T]](srcSt S, mapOp func(v T) U) Stream[U] {
	st := StreamOf[T](srcSt)
	return Stream[U]{
		srcIter: &fnIter[U]{
			fn: func(wrappedFn func(U) bool) {
				st.srcIter.iterate(func(val T) bool {
					return wrappedFn(mapOp(val))
				})
			},
		},
	}
}

// PStreamMap transforms both the values in each pair in the stream with the
// provided function.
func PStreamMap[K1, V1, K2, V2 any](
	srcSt Stream[Pair[K1, V1]],
	mapOp func(K1, V1) (K2, V2),
) Stream[Pair[K2, V2]] {
	// ques [bs]: how convinced am I that this is the right interface for this
	// behavior? particularly - could add function to map value or map key,
	// possibly in addition to or as a replacement for this.
	return StreamMap(srcSt, func(p Pair[K1, V1]) Pair[K2, V2] {
		return PairOf(mapOp(p.Get()))
	})
}

// PStreamMapValue transforms the value in each pair in the stream with the
// provided function.
func PStreamMapValue[K, V1, V2 any](
	srcSt Stream[Pair[K, V1]],
	mapOp func(K, V1) V2,
) Stream[Pair[K, V2]] {
	return StreamMap(srcSt, func(p Pair[K, V1]) Pair[K, V2] {
		return PairOf(p.First, mapOp(p.Get()))
	})
}

// PStreamMapKey transforms the key in each pair in the stream with the provided
// function.
func PStreamMapKey[K1, K2, V any](
	srcSt Stream[Pair[K1, V]],
	mapOp func(K1, V) K2,
) Stream[Pair[K2, V]] {
	return StreamMap(srcSt, func(p Pair[K1, V]) Pair[K2, V] {
		return PairOf(mapOp(p.Get()), p.Second)
	})
}

// StreamPeek will call the function on each element in the stream, but without
// any other side effects on the stream.
func StreamPeek[T any](srcSt Stream[T], peekOp func(v T)) Stream[T] {

	// so - I still feel like I'm not _quite_ consistently conceptualizing this
	// correctly. Or rather, I didn't go far enough in helpers.
	//
	// Particularly - I often do now need to just wrap an existing fn in another.
	// It's not clear to me if you really ought be interacting w/ the iteration. Seems
	// like there's some types of
	//
	// I'll also add that I'm a little worried there might be a few invisible barriers
	// in here w/ functions that as they pile up certain optimizations and inlining
	// may go out the window.
	//
	// That might start getting a bit hairy though. Roughly speaking: again, there is
	// sort of a layer different between "function that "

	return Stream[T]{
		srcIter: &fnIter[T]{
			fn: func(f func(T) bool) {
				srcSt.srcIter.iterate(func(val T) (advance bool) {
					peekOp(val)
					return f(val)
				})
			},
		},
	}
}

// StreamPeek will call the function on each pair in the stream, but without any
// other side effects on the stream.
func PStreamPeek[T, U any](
	srcSt Stream[Pair[T, U]],
	peekOp func(v T, u U),
) Stream[Pair[T, U]] {
	return StreamPeek(srcSt, func(p Pair[T, U]) {
		peekOp(p.Get())
	})
}

// StreamKeep returns a stream consisting of all elements of the source stream
// that match the given check.
func StreamKeep[T any](srcSt Stream[T], keepOp func(T) bool) Stream[T] {
	return Stream[T]{
		srcIter: &fnIter[T]{
			fn: func(f func(T) bool) {
				srcSt.srcIter.iterate(func(val T) (advance bool) {
					if keepOp(val) {
						return f(val)
					}
					return true
				})
			},
		},
	}

	// return StreamOfFn(func() Opt[T] {
	// 	for {
	// 		next := src.src.Next()
	// 		if next.IsEmpty() || check(next.UnsafeGet()) {
	// 			return next
	// 		}
	// 	}
	// })
}

// PStreamKeep returns a stream consisting of all elements of the source stream
// that match the given check.
func PStreamKeep[T, U any](
	srcSt Stream[Pair[T, U]],
	keepOp func(T, U) bool,
) Stream[Pair[T, U]] {

	return Stream[Pair[T, U]]{
		srcIter: &fnIter[Pair[T, U]]{
			fn: func(f func(Pair[T, U]) bool) {
				srcSt.srcIter.iterate(func(val Pair[T, U]) (advance bool) {
					if keepOp(val.First, val.Second) {
						return f(val)
					}
					return true
				})
			},
		},
	}
}

// StreamRemove returns a stream consisting of all elements of the source stream
// that do _not_ match the given check.
func StreamRemove[T any](srcSt Stream[T], removeOp func(T) bool) Stream[T] {
	return Stream[T]{
		srcIter: &fnIter[T]{
			fn: func(f func(T) bool) {
				srcSt.srcIter.iterate(func(val T) (advance bool) {
					if !removeOp(val) {
						return f(val)
					}
					return true
				})
			},
		},
	}
}

// PStreamRemove returns a stream consisting of all elements of the source stream
// that do _not_ match the given check.
func PStreamRemove[T, U any](
	srcSt Stream[Pair[T, U]],
	removeOp func(T, U) bool,
) Stream[Pair[T, U]] {

	return Stream[Pair[T, U]]{
		srcIter: &fnIter[Pair[T, U]]{
			fn: func(f func(Pair[T, U]) bool) {
				srcSt.srcIter.iterate(func(val Pair[T, U]) (advance bool) {
					if !removeOp(val.First, val.Second) {
						return f(val)
					}
					return true
				})
			},
		},
	}
}

// Each will perform the given function on each element of the input.
//
// Note this takes any streamable value as input - e.g. a stream or list can
// be given. Examples -
//
//    Each(Slice(1, 2, 3), func(val int) {
//       fmt.Println("value is - ", val)
//    })
//
//    Each(StreamOfVals(1, 2, 3), func(val int) {
//       fmt.Println("value is - ", val)
//    })
//
func Each[T any, S Streamable[T]](srcSt S, eachOp func(T)) {
	StreamOf[T](srcSt).Each(eachOp)
}

// EachPair will perform the given function on each element of the input pair
// stream.
func EachPair[T, U any, S Streamable[Pair[T, U]]](srcSt S, eachOp func(T, U)) {
	// ques [bs]: this is one of an increasing number of cases where restricting
	// pair to have a comparable would be really convenient.
	Each(srcSt, func(p Pair[T, U]) {
		eachOp(p.First, p.Second)
	})
}

// StreamReduce combines all the values in the stream down to one of type `U`. A
// value of `U` type is initialized, and merge is called repeatedly on each
// element in the stream with the current value of `U`. Once the stream is
// finished, the final value is returned.
//
// Example:
//
//   st := StreamOfVals(1, 2, 3)
//   sum := StreamReduce(st, func(v1, v2 int) int { return v1 + v2 })
//
func StreamReduce[T any, U any](
	srcSt Stream[T],
	reduceOp func(total U, val T) U,
) U {
	// todo [bs]: workshop the initialization here
	u := new(U)
	return StreamReduceInit(srcSt, *u, reduceOp)
}

// StreamReduceInit combines all the values in the stream down to one of type
// `U`. `merge`` is called repeatedly on each element in the stream with the
// current value of `U`, starting with the provided `initVal`. Once the stream
// is finished, the final value is returned.
//
// Example:
//
//   st := StreamOfVals(1, 2, 3)
//   product := StreamReduceInit(st, 1, func(v1, v2 int) int { return v1 * v2 })
//
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
//   st := StreamOfMap(map[string]int{
//     "a": 1,
//     "b": 2,
//     "c": 3,
//   })
//   sum := StreamReduce(st, func(total int, k string, v int) int { return total + v })
//
func PStreamReduce[T any, U any, V any](
	srcSt Stream[Pair[T, U]],
	reduceOp func(total V, first T, second U) V,
) V {
	v := new(V)
	return PStreamReduceInit(srcSt, *v, reduceOp)
}

// PStreamReduceInit combines all the values in the pair-stream down to one of
// type `V`. `merge`` is called repeatedly on each element in the stream with
// the current value of `V`, starting with the provided `initVal`. Once the
// stream is finished, the final value is returned.
//
// Example:
//
//   st := StreamOfMap(map[string]int{
//     "a": 1,
//     "b": 2,
//     "c": 3,
//   })
//   product := PStreamReduceInit(st, 1, func(total int, key string, val int) int {
//     return total * val
//   })
//
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

// todo [bs]: let's add a few simple things like find, any, first, etc.
