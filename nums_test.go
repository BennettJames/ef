package ef

import (
	"fmt"
	"math"
	"testing"
	"unsafe"
)

func Test_min_max_values(t *testing.T) {
	type S struct {
		v1 *string
	}

	fmt.Printf("@@@ uint max - %v\n", float32(math.Inf(-1)))

	s := &S{
		v1: Ptr(Deref(Ptr("hello, world!"))),
	}
	fmt.Println("val is - ", s)

	var x int = 0
	fmt.Println("sizeof int -", unsafe.Sizeof(x))
	var y any = x
	fmt.Println("sizeof any int -", unsafe.Sizeof(y))

	type A struct {
		a, b, c, d int
	}

	var a A
	fmt.Println("sizeof A -", unsafe.Sizeof(a))
	var b any = a
	fmt.Println("sizeof any A -", unsafe.Sizeof(b))

}
