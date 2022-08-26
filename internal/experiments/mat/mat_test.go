package mat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Misc(t *testing.T) {
	// quick note: pretty sure there's plenty of bugs in the current impl. Also, I
	// think I'm gonna want to add proper direct sizing info sooner rather than
	// later. Let's iterate through setting up better constructor / print / test /
	// usability stuff, with the expectation I'll have to rework a lot.
	//
	// Also - need to improve naming schemes for all this.
	//
	// ques - any vectorization schemes worthwhile here?

	t.Run("scratch", func(t *testing.T) {
		t.Skip()
		m := Zero[S2, S5]()
		Print(m)

		m.Set(0, 2, 5)
		Print(m)
	})

	t.Run("mult", func(t *testing.T) {
		// todo [bs]: obviously need better primitives here for
		// dealing w/ initialization, but just brute-forcing it for now.

		m1 := Zero[S2, S3]()
		m1.Set(0, 0, 1)
		m1.Set(0, 1, 1)
		m1.Set(0, 2, 1)
		m1.Set(1, 0, 1)
		m1.Set(1, 1, 1)
		m1.Set(1, 2, 1)

		m2 := Zero[S3, S4]()
		m2.Set(0, 0, 1)
		m2.Set(0, 1, 2)
		m2.Set(0, 2, 3)
		m2.Set(0, 3, 4)
		m2.Set(1, 0, 0)
		m2.Set(1, 1, 0)
		m2.Set(1, 2, 0)
		m2.Set(1, 3, 0)
		m2.Set(2, 0, 0)
		m2.Set(2, 1, 0)
		m2.Set(2, 2, 0)
		m2.Set(2, 3, 0)

		m3 := Mult(m1, m2)
		Print(m3)
	})

}

func Test_Misc_Mult2(t *testing.T) {
	m1 := Zero[S2, S3]()
	m1.Set(0, 0, 1)
	m1.Set(0, 1, 2)
	m1.Set(0, 2, 3)
	m1.Set(1, 0, 4)
	m1.Set(1, 1, 5)
	m1.Set(1, 2, 6)

	m2 := Zero[S3, S2]()
	m2.Set(0, 0, 7)
	m2.Set(0, 1, 8)
	m2.Set(1, 0, 9)
	m2.Set(1, 1, 10)
	m2.Set(2, 0, 11)
	m2.Set(2, 1, 12)

	expected := Zero[S2, S2]()
	expected.Set(0, 0, 58)
	expected.Set(0, 1, 64)
	expected.Set(1, 0, 139)
	expected.Set(1, 1, 154)

	m3 := Mult(m1, m2)
	Print(m3)

	assert.Equal(t, expected, m3)
}
