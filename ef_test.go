package ef

import (
	"bytes"
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

	fmt.Printf("@@@ pair is - %s\n", PairOf("v1", "v3"))

	// hmm, so this worked quite well. Need to think on that.
	v2 := ResOfPtr(ReadJSON2[A](in)).Val()
	fmt.Printf("@@@ second try - %+v\n", v2)

	// so I think this is running afoul of the null
	var in3 *bytes.Buffer
	ReadJSON[A](AutoReader(in3.Read))

}

func TestIntRange(t *testing.T) {

	// note [bs]: not sure the default should be inclusive.

	StreamEach(Range(5, 10), func(val int) {
		fmt.Printf("@@@ val - %v\n", val)
	})

	StreamEach(Range(uint64(0), uint64(3)), func(val uint64) {
		fmt.Printf("@@@ uint val - %v\n", val)
	})

	baseRange := Range(order(3, 0))
	asStrings := StreamMap(baseRange, func(v int) string {
		return fmt.Sprintf("<value as string: %v>", v)
	})
	withPeek := StreamPeek(asStrings, func(s string) {
		fmt.Printf("@@@ taking a peek at - %v\n", s)
	})
	StreamEach(withPeek, func(val string) {
		fmt.Printf("@@@ output is - %v\n", val)
	})

	// I think it might be time to do some collecting.
}
