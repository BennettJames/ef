package ef

import (
	"fmt"
)

type (
	// Void is a non-instantiable interface. Mostly useful when specifying type
	// parameters when one type is never used - e.g. you have to shoehorn in a
	// `Res[T]` type where T is never used, `Res[Void]` can be.
	Void interface {
		neverImplented()
	}
)

// Ptr wraps the provided value as a . Mostly useful for primitives in contexts
// where you'd otherwise have to declare an extra variable.
//
// Example:
//
//     // without Ptr:
//     value := "a string value"
//     fnThatTakesAStringPointer(&value)
//
//     // with Ptr:
//     fnThatTakesAStringPointer(ef.Ptr("a string value"))
//
func Ptr[V any](val V) *V {
	return &val
}

// DeRef does a "safe dereferencing" of a pointer. If the pointer points to a
// value, the value is returned; if it is null, it returns a zero-value for the
// underlying type.
//
// Example:
//
//     ef.DeRef(nil)             // == ""
//     ef.DeRef(ef.Ptr("hello")) // == "hello"
//
func DeRef[V any](val *V) V {
	if val == nil {
		return *new(V)
	}
	return *val
}

// AsType applies a type assertion to the given value, and panics if it fails.
func AsType[T any](val any) T {
	// ques [bs]: should there be a similar function for result and / or optional?
	// Yeah, probably. I'd add this isn't really super useful as-is - just plain type
	// asserts achieve the same thing. I think it can make sense as part of a broader
	// ecosystem of asserts and mappings within a space, but not standalone.
	asT, isT := val.(T)
	if !isT {
		// todo [bs]: custom error type
		panic(fmt.Errorf("bad type - expected '%T', got '%T'", new(T), val))
	}
	return asT
}

// Slice returns a slice of the given values. This is a bit of a gimmick -
// basically type inference lets you skip specifying types in some cases,
// which is of dubious usefulness.
//
// Example:
//
//   _ = []string{"a", "b", "c"}
//   _ = Slice("a", "b", "c")
//
func Slice[T any](vals ...T) []T {
	if len(vals) == 0 {
		return make([]T, 0)
	}
	return vals
}
