package opt

import "github.com/BennettJames/ef"

// Of returns an optional that has the given value stored in it.
//
// Usage note: if a pointer type is passed to this, then this will still
// be considered a "value optional" that is not empty. To convert a nil-able
// pointer to an optional, use `OptOfPtr`.
//
// Example:
//
//	var strPtr *string
//
//	opt1 := Of(opt1)
//	opt1.IsEmpty()         // false
//
//	opt2 := OptOfPtr(opt2)
//	opt2.IsEmpty()         // true
func Of[T any](val T) ef.Opt[T] {
	return ef.NewOptValue(val)
}

// Empty returns an empty optional.
func Empty[T any]() ef.Opt[T] {
	return ef.Opt[T]{}
}

// OfPtr converts a pointer to an unboxed optional. If the value is nil, then
// the optional is empty; if the value is present, then the optional contains
// it.
func OfPtr[T any](val *T) ef.Opt[T] {
	if val == nil {
		return ef.Opt[T]{}
	}
	return Of(*val)
}

// OfOk returns an empty optional if `ok` is false, and a value optional
// containing `val` if it is true.
//
// This is intended for cases where a function returns a boolean flag to
// indicate if some operation succeeqded, and returned a value in the first
// argument.
//
// Example:
//
//	matchOpt := OfOk(path.Match("pattern", pathName))
func OfOk[T any](val T, ok bool) ef.Opt[T] {
	if !ok {
		return ef.Opt[T]{}
	}
	return Of(val)
}

// MapGet looks up the key in the given map, returns an empty optional if the
// key is missing, and an optional containing the value if the key is in the map.
func MapGet[T comparable, U any](m map[T]U, key T) ef.Opt[U] {
	// note [bs]: not sure if I like the labelling for this - sorta conflicts with
	// MapFlatten.
	val, ok := m[key]
	return OfOk(val, ok)
}

// SliceGet returns an empty optional if the index is outside the bounds of
// the slice, or an optional containing the value at the index if it is in
// bounds.
func SliceGet[T any](s []T, index int) ef.Opt[T] {
	if index < 0 || index >= len(s) {
		return Empty[T]()
	}
	return Of(s[index])
}

// Map will call the provided function with any value the optional has, and
// returns a new optional with the returned value (or an empty optional if the
// original option is empty).
func Map[T any, U any](o ef.Opt[T], fn func(v T) U) ef.Opt[U] {
	if o.IsEmpty() {
		return ef.Opt[U]{}
	}
	return ef.NewOptValue(fn(o.UnsafeGet()))
}

// FlatMap calls the provided function with any value the optional has, but
// expects an optional to be returned.
func FlatMap[T any, U any](o ef.Opt[T], fn func(v T) ef.Opt[U]) ef.Opt[U] {
	if o.IsEmpty() {
		return ef.Opt[U]{}
	}
	return fn(o.UnsafeGet())
}

// Flatten reduces a nested optional down to one. If either the inner or outer
// optional is empty, then an empty optional is returned; otherwise an optional with
// the value is returned.
func Flatten[T any](o ef.Opt[ef.Opt[T]]) ef.Opt[T] {
	if o.IsEmpty() {
		return ef.Opt[T]{}
	}
	return o.UnsafeGet()
}

// Deref converts an optional of a pointer into a single optional of a
// value. Optionals containing pointers are rather odd - they can have a value
// of nil while still being nonempty. This converts a value-optional containing
// a nil pointer into a simple empty-optional for the dereferenced type.
//
// Example:
//
//	strVal := "hello"
//	strValPtr := &strVal
//	var strNilPtr *string
//
//	optWithPtrValue := OptOf(strValPtr)
//	optWithNilValue := OptOf(strNilPtr)
//
//	optWithPtrValue.IsEmpty() // false
//	optWithNilValue.IsEmpty() // false
//
//	var flatOptVal Opt[string] = Deref(optWithPtrValue)
//	var flatOptNil Opt[string] = Deref(optWithNilValue)
//
//	flatOptVal.IsEmpty() // false
//	flatOptNil.IsEmpty() // true
//
// Note that in this particular example this could have been avoided by using
// `OfPtr`, but some compositional cases will still lead to options around a
// pointer type where this can come in handy.
func Deref[T any](o ef.Opt[*T]) ef.Opt[T] {
	if o.IsEmpty() {
		return ef.Opt[T]{}
	}
	return OfPtr(o.UnsafeGet())
}
