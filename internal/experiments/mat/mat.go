package mat

import (
	"fmt"
	"strconv"

	"github.com/BennettJames/ef"
)

type (
	S0 struct {
	}

	S1 struct {
	}

	S2 struct {
	}

	S3 struct {
	}

	S4 struct {
	}

	S5 struct {
	}

	Size interface {
		S0 | S1 | S2 | S3 | S4 | S5
	}

	Matrix[R, C Size] struct {
		rows, cols int
		vals       []float64
	}

	matrixImpl[R, C Size, N ef.AllNumber] struct {
		// todo [bs]: let's experiment w/ this. I'm particularly interested
		// in complex support.
		//
		// also - would it possibly make sense to experiment w/ row vs column
		// major? Maybe, but I care less about that.
		rows, cols int
		vals       []N
	}

	Vector[S Size] struct {
		vals []float64
	}
)

// sidenote [bs]: this particular experiment feels like it'd be more appropriate
// for my typescript graphics project - it would actually have a use, and the
// type system there is more powerful in a way that might actually make it vaguely
// useful.

func Zero[R, C Size]() Matrix[R, C] {
	numRows, numCols := sizeToInt[R](), sizeToInt[C]()
	vals := make([]float64, numRows*numCols)
	return Matrix[R, C]{
		rows: numRows,
		cols: numCols,
		vals: vals,
	}
}

func (m *Matrix[R, C]) index(row, col int) int {
	// ques [bs]: should there be any guards here? I'd say not at this
	// level - it might not be a bad idea on higher-level methods, but
	// I think that should check the top-level values rather than
	return row*m.cols + col
}

func (m *Matrix[R, C]) Set(row, col int, val float64) {

	m.vals[m.index(row, col)] = val
}

func (m *Matrix[R, C]) Get(row, col int) float64 {
	return m.vals[m.index(row, col)]
}

func Add[R, C Size](m1, m2 Matrix[R, C]) Matrix[R, C] {
	vals := make([]float64, m1.rows*m1.cols)
	for rI := 0; rI < m1.rows; rI++ {
		for cI := 0; cI < m1.cols; cI++ {
			index := m1.index(rI, cI)
			vals[index] = m1.vals[index] + m2.vals[index]
		}
	}
	return Matrix[R, C]{
		vals: vals,
	}
}

func Mult[R1, C1, C2 Size](m1 Matrix[R1, C1], m2 Matrix[C1, C2]) Matrix[R1, C2] {

	out := Zero[R1, C2]()

	for outCol := 0; outCol < out.cols; outCol++ {
		for outRow := 0; outRow < out.rows; outRow++ {
			subTotal := float64(0)
			for m2Row := 0; m2Row < m2.rows; m2Row++ {
				subTotal += m1.Get(outRow, m2Row) * m2.Get(m2Row, outCol)
			}
			out.Set(outRow, outCol, subTotal)
		}
	}

	return out
}

func Print[R1, C1 Size](m Matrix[R1, C1]) {
	fmt.Println(Sprint(m))
}

func Sprint[R1, C1 Size](m Matrix[R1, C1]) string {
	// let's start w/ a string literal, then go from there.

	out := "[\n"
	for row := 0; row < m.rows; row++ {
		out += "  [ "

		for col := 0; col < m.cols; col++ {
			// note [bs]: probably want a precision control / limiter;
			// plus some general column-izing formatting. Maybe scientific
			// notation in some cases?
			//
			// For now, let's try for some sane-ish defaults that'll lead
			// to a somewhat functional debugging experience.
			out += strconv.FormatFloat(m.Get(row, col), 'g', 3, 64)
			out += " "
		}

		out += "]\n"
	}

	out += "]"

	return out
}

func sizeToInt[S Size]() int {
	switch any(*new(S)).(type) {
	case S0:
		return 0
	case S1:
		return 1
	case S2:
		return 2
	case S3:
		return 3
	case S4:
		return 4
	case S5:
		return 5
	default:
		panic("unreachable")
	}
}
