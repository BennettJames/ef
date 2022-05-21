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
