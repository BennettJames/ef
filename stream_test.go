package ef

import (
	"fmt"
	"testing"
)

func TestStreamStats(t *testing.T) {
	// I find some of the casting behavior here a bit suprising. Both that it
	// appears to let me only specify one type arg for streamOf, and that I had to
	// at all. Let's see if I can understand why.
	//
	// - First, just want a sense if maybe there are times a second type argument
	// can be unspecified because it's obvious when the first is not.
	//
	// - also want a sense for why the type arg is even needed. that's going to
	// be harder. I suspect that inference is probably complicated, and hard to
	// guess when you push limits.
	ary := []int{5, 10, 9, 8, 22}
	s := StreamOf[int](ary)
	stats := StreamStats(s)
	fmt.Printf("got stats - %+v\n", stats)

	// so, I think this clarifies the first point - trailing type args can be
	// elided if
	//
	// that's good to know - that might actually simplify the json thing you were
	// trying to do.
	var _ int32 = inferTest[int32, string]("hello")
}

func inferTest[N Number, T any, S any](s S) N {
	return MaxNumber[N]()
}
