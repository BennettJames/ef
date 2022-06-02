package ef

type (
	// Opt represents an optional of the given type. An optional can either have a value,
	// or be empty.
	//
	// Used well, this works as a kind of typesafe pointer - a series of helper functions
	// make it easier to compose around and use in a safe way, and discourage any behavior
	// that can result in a nil exception.
	Opt[T any] struct {
		value   T
		present bool
	}
)

// OptOf returns an optional that has the given value stored in it.
//
// Usage note: if a pointer type is passed to this, then this will still
// be considered a "value optional" that is not empty. To convert a nil-able
// pointer to an optional, use `OptOfPtr`.
//
// Example:
//
//   var strPtr *string
//
//   opt1 := OptOf(opt1)
//   opt1.IsEmpty()         // false
//
//   opt2 := OptOfPtr(opt2)
//   opt2.IsEmpty()         // true
//
func OptOf[T any](val T) Opt[T] {
	return Opt[T]{
		value:   val,
		present: true,
	}
}

// OptEmpty returns an empty optional.
func OptEmpty[T any]() Opt[T] {
	return Opt[T]{}
}

// OptOfPtr converts a pointer to an unboxed optional. If the value is nil, then
// the optional is empty; if the value is present, then the optional contains
// it.
func OptOfPtr[T any](val *T) Opt[T] {
	if val == nil {
		return Opt[T]{}
	}
	return OptOf(*val)
}

// OptOfOk returns an empty optional if `ok` is false, and a value optional
// containing `val` if it is true.
//
// This is intended for cases where a function returns a boolean flag to
// indicate if some operation succeeqded, and returned a value in the first
// argument.
//
// Example:
//
//   matchOpt := OptOfOk(path.Match("pattern", pathName))
//
func OptOfOk[T any](val T, ok bool) Opt[T] {
	if !ok {
		return Opt[T]{}
	}
	return OptOf(val)
}

// OptMapGet looks up the key in the given map, returns an empty optional if the
// key is missing, and an optional containing the value if the key is in the map.
func OptMapGet[T comparable, U any](m map[T]U, key T) Opt[U] {
	// note [bs]: not sure if I like the labelling for this - sorta conflicts with
	// MapFlatten.
	val, ok := m[key]
	return OptOfOk(val, ok)
}

// OptSliceGet returns an empty optional if the index is outside the bounds of
// the slice, or an optional containing the value at the index if it is in
// bounds.
func OptSliceGet[T any](s []T, index int) Opt[T] {
	if index < 0 || index >= len(s) {
		return OptEmpty[T]()
	}
	return OptOf(s[index])
}

// UnsafeGet returns the value if present, and panics if it does not exist. Note
// this is a dangerous method to use - generally it's best to an alternative to
// safely process, like `IfVal`, `Or`, or `OptMap`. Aim to structure the usage
// of the optional so the code can't err and makes no assumptions.
func (o Opt[T]) UnsafeGet() T {
	if !o.present {
		// todo [bs]: I still sorta suspect that the nil error should be able
		// to contain and communicate some amount of context.
		panic(&UnexpectedNilError{})
	}
	return o.value
}

// GetPtr returns a pointer of the inner value (and is nil if the optional
// is empty).
func (o Opt[T]) GetPtr() *T {
	if !o.present {
		return nil
	}
	return &o.value
}

// HasVal indicates if the optional has a value.
func (o Opt[T]) HasVal() bool {
	return o.present
}

// IsEmpty indicates if the optional lacks a value.
func (o Opt[T]) IsEmpty() bool {
	return !o.present
}

// todo [bs]: add a matching IsEmpty pure function once this is moved to
// a subpackage for filtering.

// IfVal executes the provided function with the stored value if the optional
// has a value; otherwise does nothing. Returns itself for chaining.
func (o Opt[T]) IfVal(fn func(v T)) Opt[T] {
	if o.present {
		fn(o.value)
	}
	return o
}

// IfEmpty calls the passed function if the optional is empty, otherwise does
// nothing. Returns itself for chaining.
func (o Opt[T]) IfEmpty(fn func()) Opt[T] {
	if !o.present {
		fn()
	}
	return o
}

// Or returns the provided value if the optional is empty, or the value if it
// has one.
func (o Opt[T]) Or(altVal T) T {
	if o.HasVal() {
		return o.value
	}
	return altVal
}

// OrCalc calls and returns the value from the function if the optional is
// empty, the the value if the optional has one.
//
// This is an alternative to `Or` for cases where it may be undesirable to
// unnecessarily compute the alternative value - for instance, if the
// calculation is expensive.
func (o Opt[T]) OrCalc(fn func() T) T {
	if o.HasVal() {
		return o.value
	}
	return fn()
}

// ToList converts the optional to a list. If the optional is empty, then the
// list is; otherwise it consists of just the single value held by the optional.
func (o Opt[T]) ToList() []T {
	if o.HasVal() {
		return []T{o.value}
	} else {
		return []T{}
	}
}

// OptMap will call the provided function with any value the optional has, and
// returns a new optional with the returned value (or an empty optional if the
// original option is empty).
func OptMap[T any, U any](o Opt[T], fn func(v T) U) Opt[U] {
	if !o.present {
		return Opt[U]{}
	}
	return OptOf(fn(o.value))
}

// OptFlatMap calls the provided function with any value the optional has, but
// expects an optional to be returned.
func OptFlatMap[T any, U any](o Opt[T], fn func(v T) Opt[U]) Opt[U] {
	if !o.present {
		return Opt[U]{}
	}
	return fn(o.value)
}

// OptFlatten reduces a nested optional down to one. If either the inner or outer
// optional is empty, then an empty optional is returned; otherwise an optional with
// the value is returned.
func OptFlatten[T any](o Opt[Opt[T]]) Opt[T] {
	if !o.present {
		return Opt[T]{}
	}
	return o.value
}

// OptFlattenPtr converts an optional of a pointer into a single optional of a
// value. Optionals containing pointers are rather odd - they can have a value
// of nil while still being nonempty. This converts a value-optional containing
// a nil pointer into a simple empty-optional for the dereferenced type.
//
// Example:
//
//   strVal := "hello"
//   strValPtr := &strVal
//   var strNilPtr *string
//
//   optWithPtrValue := OptOf(strValPtr)
//   optWithNilValue := OptOf(strNilPtr)
//
//   optWithPtrValue.IsEmpty() // false
//   optWithNilValue.IsEmpty() // false
//
//   var flatOptVal Opt[string] = OptFlattenPtr(optWithPtrValue)
//   var flatOptNil Opt[string] = OptFlattenPtr(optWithNilValue)
//
//   flatOptVal.IsEmpty() // false
//   flatOptNil.IsEmpty() // true
//
// Note that in this particular example this could have been avoided by using
// `OptOfPtr`, but some compositional cases will still lead to options around a
// pointer type where this can come in handy.
func OptFlattenPtr[T any](o Opt[*T]) Opt[T] {
	if !o.present {
		return Opt[T]{}
	}
	return OptOfPtr(o.value)
}
