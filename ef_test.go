package ef

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {

}

func TestScratch(t *testing.T) {
	var iter Iterator[int] = &listIter[int]{
		list:      []int{1, 2, 3},
		nextIndex: 0,
	}

	IteratorForEach(iter, func(v int) {
		fmt.Println("@@@ value is - ", v)
	})
}

func TestReadJSON(t *testing.T) {
	type A struct {
		Foo string
		Bar string
	}
	expected := A{"foo", "bar"}

	// so - not at all convinced this is like _better_, but it's interesting.
	//
	// Having to state the type largely eliminates any advantage this might hope
	// to have. There _might_ be a handful of cases where that wouldn't be true,
	// but I'd expect it to be rare.
	//
	// On the whole, I'm going to rate this "not quite powerful enough to be
	// useful", which again I think might really be the design principles in
	// action here. Nothing wrong with that.
	//
	// Not that this is a good idea, but could you elide the first argument via a
	// form of currying? I kinda doubt that'd work here, but let's play around.
	var actual *A = MustReadJSON[string, A](`{"Foo": "foo", "Bar": "bar"}`)
	fmt.Printf("equal? '%+v' '%+v'", expected, *actual)
}
