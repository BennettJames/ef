package ef

import (
	"reflect"
)

// IsEmptySlice indicates if the provided slice is empty. This is mainly intended
// for filter. Example:
//
//   st := StreamOfVals(Slice(1, 2), Slice(3), Slice(), Slice(4))
//   stWithoutEmpty := stream.Filter(st, stream.IsEmptySlice[int])
//
func IsEmptySlice[T any](slice []T) bool {
	return len(slice) == 0
}

// IsEmptyStr indicates if the provided string is empty. This is mainly intended
// for filter. Example:
//
//   st := StreamOfVals("ab", "", "c")
//   stWithoutEmpty := stream.Filter(st, stream.IsEmptyStr)
//
func IsEmptyStr(str string) bool {
	return len(str) == 0
}

// IsEmptyMap indicates if the provided map is empty. This is mainly intended
// for filter. Example:
//
//   st := StreamOfVals(map[int]string{}, map[int]string{ 1: "a" })
//   stWithoutEmpty := stream.Filter(st, stream.IsEmptyMap[int, string])
//
func IsEmptyMap[T comparable, U any](m map[T]U) bool {
	return len(m) == 0
}

func isReallyNil(i any) bool {
	// note [bs]: not yet quite sure if I want this. Let's do a bit more testing.
	return i == nil || reflect.ValueOf(i).IsNil()
}

func And[T any](checkFns ...func(T) bool) func(T) bool {
	return func(val T) bool {
		for _, fn := range checkFns {
			if !fn(val) {
				return false
			}
		}
		return true
	}
}

func Or[T any](checkFns ...func(T) bool) func(T) bool {
	return func(val T) bool {
		for _, fn := range checkFns {
			if fn(val) {
				return true
			}
		}
		return len(checkFns) == 0
	}
}

func Not[T any](checkFn func(T) bool) func(T) bool {
	return func(t T) bool {
		return !checkFn(t)
	}
}

func Equal[T comparable](checkVal T) func(T) bool {
	return func(val T) bool {
		return checkVal == val
	}
}

func Greater[T Number](checkVal T) func(T) bool {
	return func(val T) bool {
		// ques [bs] - is this ordering sensible?
		return checkVal < val
	}
}

func GreaterOrEqual[T Number](checkVal T) func(T) bool {
	return func(val T) bool {
		// ques [bs] - is this ordering sensible?
		return checkVal <= val
	}
}

func Lesser[T Number](checkVal T) func(T) bool {
	return func(val T) bool {
		// ques [bs] - is this ordering sensible?
		return checkVal > val
	}
}

func LesserOrEqual[T Number](checkVal T) func(T) bool {
	return func(val T) bool {
		// ques [bs] - is this ordering sensible?
		return checkVal >= val
	}
}
