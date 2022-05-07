package ef

import (
	"fmt"
	"math"
	"testing"
)

func Test_min_max_values(t *testing.T) {
	fmt.Printf("@@@ uint max - %v\n", float32(math.Inf(-1)))

	s := &S{
		v1: Ptr(Deref(Ptr("hello, world!"))),
	}
	fmt.Println("val is - ", s)

}

func Ptr[V any](val V) *V {
	return &val
}

func Deref[V any](val *V) V {
	return *val
}

type S struct {
	v1 *string
}
