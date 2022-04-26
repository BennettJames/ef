package ef

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {

}

func TestScratch(t *testing.T) {
	var iter Iter[int] = &listIter[int]{
		list:      []int{1, 2, 3},
		nextIndex: 0,
	}

	IterEach(iter, func(v int) {
		fmt.Println("@@@ value is - ", v)
	})
}

func TestReadJSON(t *testing.T) {
	type A struct {
		Foo string
		Bar string
	}
	// expected := A{"foo", "bar"}

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

	in := `{"Foo": "foo", "Bar": "bar"}`

	// so this is interesting. Still would be better to be able to take
	// an interface in AutoReader, but still decent.
	v, err := ReadJSON[A](AutoReader(in))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("@@@ read value - '%+v'\n", v)

	fmt.Printf("@@@ pair is - %s\n", NewPair("v1", "v3"))

}
