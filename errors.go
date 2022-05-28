package ef

import "fmt"

type (
	RecoverError struct {
		recovered any
	}

	UnexpectedNilError struct{}

	// UnreachableError is designed to be thrown
	UnreachableError struct{}
)

func (e *RecoverError) Error() string {
	// note [bs]: not super happy with this text value; let's workshop it.
	return fmt.Sprintf("Recovered try with value: '%v'", e.recovered)
}

func (e *UnexpectedNilError) Error() string {
	// note [bs]: I don't think this type and it's behavior 100% make sense as is,
	// but I feel like I might be circling towards something more meaningful.

	// ques [bs]: is there any decent way to make this more contexually
	// useful? Let's think on whether adding a bit of calling context
	// might be useful. I suspect yes.
	return "Result encountered an unexpected nil"
}

func (e *UnreachableError) Error() string {
	// ques [bs]: should I add any context-building facilities
	// to this? yeah, this feels plain to the point of uselessness.

	return "unreachable"
}

func Recover(errAddr *error) {
	if errAddr == nil {
		panic("Recover called with nil result reference")
	}

	switch narrowed := recover().(type) {
	case nil:
		// ques [bs]: is doing this in a type switch less efficient then just checking
		// directly?
		return
	case error:
		*errAddr = narrowed
	default:
		*errAddr = &RecoverError{recovered: narrowed}
	}
}

func Catch(errAddr *error, recoverFn func(error) error) {
	if errAddr == nil {
		panic("Recover called with nil result reference")
	}

	var recoveredErr error
	switch narrowed := recover().(type) {
	case nil:
		// ques [bs]: is doing this in a type switch less efficient then just checking
		// directly?
		return
	case error:
		recoveredErr = narrowed
	default:
		recoveredErr = &RecoverError{recovered: narrowed}
	}
	*errAddr = recoverFn(recoveredErr)
}

func Try[T any](v T, err error) T {
	if err != nil {
		// todo [bs]: probably should wrap this and unwrap in a recover
		panic(err)
	}
	return v
}

func Try2[T, U any](t T, u U, err error) (T, U) {
	if err != nil {
		// todo [bs]: probably should wrap this and unwrap in a recover
		panic(err)
	}
	return t, u
}

func Try3[T, U, V any](t T, u U, v V, err error) (T, U, V) {
	if err != nil {
		// todo [bs]: probably should wrap this and unwrap in a recover
		panic(err)
	}
	return t, u, v
}

func Assert(check bool, msg string) {
	if !check {
		// todo [bs]: wrap this
		panic(msg)
	}
}

func Assertf(check bool, msg string, a ...any) {
	if !check {
		// todo [bs]: wrap this better
		panic(fmt.Errorf(msg, a...))
	}
}
