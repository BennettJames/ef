package ef

import (
	"fmt"
	"math"
	"testing"
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

}
