package ef

type (
	Opt[T any] struct {
		value   T
		present bool
	}

	// OptLike is a union of either an optional or a pointer. This is mostly
	// for a few convenience functions, where being able to take either is
	// useful.
	OptLike[T any] interface {
		Opt[T] | ~*T
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

// OptOfPtr converts a pointer to an unboxed optional. If the value
// is nil, then the optional is empty; if the value is present, then
// the optional contains it.
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
// indicate if some operation succeeded, and returned a value in the first
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

// Get returns the value if present, and panics if it does not exist. Note this
// is a dangerous method to use - generally it's best to an alternative to
// safely process, like `IfVal`, `Or`, or `OptMap`. Aim to structure the usage
// of the optional so the code can't err and makes no assumptions.
func (o Opt[T]) Get() T {
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

// IsVal indicates if the optional has a value.
func (o Opt[T]) IsVal() bool {
	return o.present
}

// IsEmpty indicates if the optional lacks a value.
func (o Opt[T]) IsEmpty() bool {
	return !o.present
}

// IfVal executes the provided function with the stored value if the optional
// has a value; otherwise does nothing.
func (o Opt[T]) IfVal(fn func(v T)) {
	if o.present {
		fn(o.value)
	}
}

func (o Opt[T]) Or(altVal T) T {
	if o.IsVal() {
		return o.value
	}
	return altVal
}

func (o Opt[T]) OrCalc(orFn func() T) T {
	if o.IsVal() {
		return o.value
	}
	return orFn()
}

func (o Opt[T]) ToList() []T {
	if o.IsVal() {
		return []T{o.value}
	} else {
		return nil
	}
}

func OptMap[T any, U any](o Opt[T], fn func(v T) U) Opt[U] {
	if !o.present {
		return Opt[U]{}
	}
	return OptOf(fn(o.value))
}

func OptFlatMap[T any, U any](o Opt[T], fn func(v T) Opt[U]) Opt[U] {
	if !o.present {
		return Opt[U]{}
	}
	return fn(o.value)
}

// OptFlatten reduces a nested optional down to one.
func OptFlatten[T any, O OptLike[T]](o Opt[O]) Opt[T] {
	if !o.present {
		return OptEmpty[T]()
	}
	return narrowOptLike[T](o.value)
}

func narrowOptLike[T any, O OptLike[T]](o O) Opt[T] {
	switch narrowed := any(o).(type) {
	case Opt[T]:
		return narrowed
	case *T:
		if narrowed == nil {
			return OptEmpty[T]()
		}
		return OptOf(*narrowed)
	default:
		panic(&UnreachableError{})
	}
}
