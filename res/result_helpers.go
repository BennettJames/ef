package res

import "github.com/BennettJames/ef"

// Of creates a result from a pair of values. If error is set, it will be an
// error type; otherwise a result.
//
// Note that if both are set / nonnull, the value is not kept - just the error.
func Of[T any](val T, e error) ef.Res[T] {
	if e == nil {
		return Val(val)
	}
	return ResErr[T](e)
}

// Creates a result from three values, where the first two are turned into a
// pair with the values stored.
//
// This is mostly designed for cases where a function returns three values and
// you wish to wrap it in a result; e.g. -
//
//	var r io.RuneReader
//	var res Res[Pair[rune, int]] = result.Of2(r.ReadRune())
func Of2[T, U any](v1 T, v2 U, e error) ef.Res[ef.Pair[T, U]] {
	if e == nil {
		return Val(ef.PairOf(v1, v2))
	}
	return ResErr[ef.Pair[T, U]](e)
}

// OfPtr takes a par of a pointer value and an error, and converts it to a a
// result. If the error is nonnil, then the result is an error type with the
// error stored. If the value is present, then the pointer's value is stored in
// the result. If both are nil, then an error result with a nil pointer error is
// returned.
func OfPtr[T any](val *T, e error) ef.Res[T] {
	if e != nil {
		return ResErr[T](e)
	}
	if val != nil {
		return Val(*val)
	}
	return ResErr[T](&ef.UnexpectedNilError{})
}

// OfOpt will return an value type result if the optional has a value, or a
// result with a nil reference error.
func OfOpt[T any](o ef.Opt[T]) ef.Res[T] {
	// note [bs]: pretty sure this equivalent of ResOfPtr(o.Get())
	return ef.OptMap(o, func(v T) ef.Res[T] {
		return Val(v)
	}).OrCalc(func() ef.Res[T] {
		return ResErr[T](&ef.UnexpectedNilError{})
	})
}

// Val creates a result from the provided value.
func Val[T any](val T) ef.Res[T] {
	return ef.NewResValue(val)
}

// Deref will create a result with a value from the pointer. If the pointer is
// nil, a zero value is used.
func Deref[T any](r ef.Res[*T]) ef.Res[T] {
	// ques [bs]: how happy am I with this behavior _really_?
	return Map(r, func(val *T) T {
		return ef.DeRef(val)
	})
}

// ResErr creates an error result for the given error. If the error is null,
// then the result is a value type with a zero V value.
func ResErr[T any](err error) ef.Res[T] {
	return ef.NewResError[T](err)
}

// Map will execute the passed function if the result is a value; otherwise if
// an error result returns the error.
func Map[T, U any](r ef.Res[T], fn func(val T) U) ef.Res[U] {
	if r.IsVal() {
		return Val(fn(r.Val()))
	}
	return ResErr[U](r.Err())
}

// FlatMap is as ResMap, but expects a result from the inner function.
func FlatMap[T, U any](r ef.Res[T], fn func(val T) ef.Res[U]) ef.Res[U] {
	if r.IsVal() {
		return fn(r.Val())
	}
	return ResErr[U](r.Err())
}

// Recover performs automatic recovery from a panic, and converts the panic to
// an error result.
//
// Example -
//
//	func appendStrings(v1, v2 *string) (res Res[string]) {
//	    defer Recover(&res)
//	    if v1 == nil || v2 == nil {
//	        panic("unexpected null strings")
//	    }
//	    return ResOfVal(*v1 + *v2)
//	}
//
// The above code will return an error result if either pointer is nil.
//
// Recover must be given a valid address. It panic if the given result pointer
// is nil.
func Recover[T any](r *ef.Res[T]) {
	if r == nil {
		// todo [bs]: consider defining a custom error for this
		panic("ResRecover called with nil result reference")
	}

	switch narrowed := recover().(type) {
	case nil:
		// ques [bs]: is doing this in a type switch less efficient then just checking
		// directly?
		return
	case error:
		// todo [bs]: consider unwrapping certain internal error types here - e.g.
		// don't
		*r = ResErr[T](narrowed)
	default:
		*r = ResErr[T](ef.NewRecoverError(narrowed))
	}
}

// TryMap will execute the given function with the result value, provided the
// result is a value type. If it is an error type, then the error is returned.
//
// Any panics in the inner function will be converted to an error result.
func TryMap[V, U any](r ef.Res[V], fn func(val V) U) (res ef.Res[U]) {
	defer Recover(&res)
	if !r.IsVal() {
		return ResErr[U](r.Err())
	}
	return Val(fn(r.Val()))
}

// TryFlatMap is as Try, but expects a result from the inner function.
func TryFlatMap[V, U any](
	r ef.Res[V],
	fn func(val V) ef.Res[U],
) (res ef.Res[U]) {
	defer Recover(&res)
	if !r.IsVal() {
		return ResErr[U](r.Err())
	}
	return Flatten(Val(fn(r.Val())))
}

// Flatten turns a nested result into a single flat one - if either the inner
// or the outer result has an error, then it is returned as an error result.
// Otherwise, the inner value is returned.
//
// Examples:
//
//	Flatten(ResOfVal(ResOfVal[string]("value")))             // == ResOfVal[string]("value")
//	Flatten(ResOfVal(ResOfErr[string](fmt.Errorf("error")))) // == ResOfErr[string](fmt.Errorf("error"))
//	Flatten(ResOfErr[Res[string]](fmt.Errorf("error")))      // == ResOfErr[string](fmt.Errorf("error"))
func Flatten[T any](r ef.Res[ef.Res[T]]) ef.Res[T] {
	if r.IsErr() {
		return ResErr[T](r.Err())
	}
	return r.Val()
}
