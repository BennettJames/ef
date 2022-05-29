package ef

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestPair(t *testing.T) {

	t.Run("Of", func(t *testing.T) {
		assert.Equal(t,
			Pair[string, int]{
				First:  "hello",
				Second: 22,
			},
			PairOf("hello", 22))
	})

	t.Run("Get", func(t *testing.T) {
		v1, v2 := PairOf("hello", 22).Get()
		assert.Equal(t, "hello", v1)
		assert.Equal(t, 22, v2)
	})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, "(`hello`, `22`)", PairOf("hello", 22).String())
	})
}

func TestSize(t *testing.T) {
	type FooI interface {
		DoIt()
	}

	foo := MakeIt[FooI]()
	fooPtr := MakeItPtr[FooI]()

	fmt.Printf("@@@ foo - %T\n", foo)       // nil
	fmt.Printf("@@@ fooPtr - %T\n", fooPtr) // *ef.FooI

	fmt.Printf("@@@ foo size - %v\n", unsafe.Sizeof(foo))       // 16
	fmt.Printf("@@@ fooPtr size - %v\n", unsafe.Sizeof(fooPtr)) // 8
}

func MakeIt[T any]() T {
	return *new(T)
}

func MakeItPtr[T any]() *T {
	return new(T)
}
