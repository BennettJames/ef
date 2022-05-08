package ef

import "fmt"

// Res is a "result" that tracks a value that can be either an error
// or a value, but not both. It comes with a set of methods and function
// helpers to make it easier to correctly work with errors.
//
// A note on nil - if a result wraps a nilable type (e.g. a pointer),
// then it is possible that the value and the error can be nil. This is
// still considered a "value typed result".
type Res[T any] struct {
	val T
	err error
}

// Get returns both the value and the error of the result.
func (r Res[T]) Get() (T, error) {
	return r.val, r.err
}

// GetPtr returns the value / error pair, but with the value wrapped as a
// pointer.
func (r Res[T]) GetPtr() (*T, error) {
	if r.err != nil {
		return nil, r.err
	}
	// ques [bs]: does this expose any weird undesirable mutability?
	return &r.val, nil
}

// Val returns the underlying value from the result, or panics if the result is
// not a value type.
func (r Res[T]) Val() T {
	if !r.IsVal() {
		panic("res.Val() called on non-value result")
	}
	return r.val
}

// Err returns the underlying error from the result, or panics if the result
// is not an error type.
func (r Res[T]) Err() error {
	if r.IsVal() {
		panic("res.Err() called on non-error result")
	}
	return r.err
}

// String is just a simple string representation of the result for debugging.
func (r Res[T]) String() string {
	// ques [bs]: should this have more of a structural difference between
	// value / error?
	return fmt.Sprintf("<val='%v' err='%v'>", r.val, r.err)
}

// ResOf creates a result from a pair of values. If error is set, it will be
// an error type; otherwise a result.
//
// Note that if both are set / nonnull, the value is not kept - just the
// error.
func ResOf[T any](val T, e error) Res[T] {
	if e == nil {
		return ResVal(val)
	}
	return ResErr[T](e)
}

// ResOfPtr takes a par of a pointer value and an error, and converts it to a a
// result. If the error is nonnil, then the result is an error type with the
// error stored. If the value is present, then the pointer's value is stored in
// the result. If both are nil, then an error result with a nil pointer error is
// returned.
func ResOfPtr[T any](val *T, e error) Res[T] {
	// note [bs]: not 100% convinced this function need exist. Let's think
	// a bit about composition here.
	if e != nil {
		return ResErr[T](e)
	}
	if val != nil {
		return ResVal(*val)
	}
	return ResErr[T](&ResultNilError{})
}

func ResDeref[T any](r Res[*T]) Res[T] {
	return ResMap(r, func(val *T) T {
		return Deref(val)
	})
}

// ResVal
func ResVal[T any](val T) Res[T] {
	return Res[T]{
		val: val,
	}
}

// ResErr creates an error result for the given error. If the error is
// null, then the result is a value type with a zero V value.
func ResErr[T any](e error) Res[T] {
	return Res[T]{
		err: e,
	}
}

// IsVal indicates if the result has a value (and is not an error).
func (r Res[T]) IsVal() bool {
	return r.err == nil
}

// IsErr indicates if the result is an error (and does not have a value).
func (r Res[T]) IsErr() bool {
	return r.err != nil
}

// IfVal will execute the passed function if the result is a value.
func (r Res[T]) IfVal(fn func(val T)) {
	if r.IsVal() {
		fn(r.val)
	}
}

// IfErr will execute the passed function if the result is an error.
func (r Res[T]) IfErr(fn func(e error)) {
	if !r.IsVal() {
		fn(r.err)
	}
}

// ResMap will execute the passed function if the result is a value; otherwise if
// an error result returns the error.
func ResMap[T, U any](r Res[T], fn func(val T) U) Res[U] {
	if r.IsVal() {
		return ResVal(fn(r.val))
	}
	return ResErr[U](r.err)
}

// ResFlatMap is as ResMap, but expects a result from the inner function.
func ResFlatMap[T, U any](r Res[T], fn func(val T) Res[U]) Res[U] {
	if r.IsVal() {
		return fn(r.val)
	}
	return ResErr[U](r.err)
}

// ResRecover performs automatic recovery from a panic, and converts the panic
// to an error result.
//
// Example -
//
//   func appendStrings(v1, v2 *string) (res Res[string]) {
//       defer ResRecover(&res)
//       if v1 == nil || v2 == nil {
//           panic("unexpected null strings")
//       }
//       return ResOfVal(*v1 + *v2)
//   }
//
// The above code will return an error result if either pointer is nil.
//
// ResRecover must be given a valid address. It panic if the given result
// pointer is nil.
func ResRecover[T any](r *Res[T]) {
	if r == nil {
		// todo [bs]: should define a custom error for this
		panic("ResRecover called with nil result reference")
	}
	if rec := recover(); rec != nil {
		if err, isErr := rec.(error); isErr {
			*r = ResErr[T](err)
		} else {
			*r = ResErr[T](&ResultRecoverError{recovered: rec})
		}
	}
}

// ResTry will execute the given function with the result value, provided the
// result is a value type. If it is an error type, then the error is returned.
//
// Any panics in the inner function will be converted to an error result.
func ResTry[V, U any](r Res[V], fn func(val V) U) (res Res[U]) {
	defer ResRecover(&res)
	if !r.IsVal() {
		return ResErr[U](r.err)
	}
	return ResVal(fn(r.val))
}

// ResFlatTry is as ResTry, but expects a result from the inner function.
func ResFlatTry[V, U any](r Res[V], fn func(val V) Res[U]) (res Res[U]) {
	defer ResRecover(&res)
	if !r.IsVal() {
		return ResErr[U](r.err)
	}
	return ResFlatten(ResVal(fn(r.val)))
}

// todo [bs]: consider use cases for returning results from try / map. Could
// just let the user all flatten on them, or could make custom "flat" calls
// for those. I'd lean towards the later - just a bit cleaner.

// ResFlatten turns a nested result into a single flat one - if either the inner
// or the outer result has an error, then it is returned as an error result.
// Otherwise, the inner value is returned.
//
// Examples:
//
//    ResFlatten(ResOfVal(ResOfVal[string]("value")))             // == ResOfVal[string]("value")
//    ResFlatten(ResOfVal(ResOfErr[string](fmt.Errorf("error")))) // == ResOfErr[string](fmt.Errorf("error"))
//    ResFlatten(ResOfErr[Res[string]](fmt.Errorf("error")))      // == ResOfErr[string](fmt.Errorf("error"))
//
func ResFlatten[T any](r Res[Res[T]]) Res[T] {
	if r.IsErr() {
		return ResErr[T](r.err)
	}
	return r.val
}

type ResultRecoverError struct {
	recovered any
}

func (e *ResultRecoverError) Error() string {
	// note [bs]: not super happy with this text value; let's workshop it.
	return fmt.Sprintf("Recovered try with value: '%v'", e.recovered)
}

type ResultNilError struct{}

func (e *ResultNilError) Error() string {
	// note [bs]: I don't think this type and it's behavior 100% make sense as is,
	// but I feel like I might be circling towards something more meaningful.
	return "Result encountered an unexpected nil"
}
