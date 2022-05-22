package rangebench2

import "testing"

func BenchmarkRangeStructVsClosure(b *testing.B) {

	// Result:
	//
	// goos: darwin
	// goarch: amd64
	// pkg: github.com/bennettjames/ef/internal/bench/rangebench
	// cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
	// BenchmarkRangeStructVsClosure
	// BenchmarkRangeStructVsClosure/struct
	// BenchmarkRangeStructVsClosure/struct-16         	  175710	      6512 ns/op	   25232 B/op	      13 allocs/op
	// BenchmarkRangeStructVsClosure/closure
	// BenchmarkRangeStructVsClosure/closure-16        	  163255	      7706 ns/op	   25272 B/op	      15 allocs/op
	// PASS
	// ok  	github.com/bennettjames/ef/internal/bench/rangebench	2.645s

	// so, interesting thing to point out - the "native test" has a wider gap -
	// 7100 ns / op vs 9200, respectively.
	//
	// So, obvious question - why? Possible this is just a matter of the functions you
	// inlined. While they were minor functions I honestly kinda expected to be inlined anyway,
	// possible they were "real" in some sense.

	// alright, so tried to undo the inlining. Now getting about the same - ~82%
	// for the range struct, vs 76% on the original (also the base one has worse
	// overall performance, but not sure that's meaningful?)
	//
	// I wouldn't mind digging into that a bit further. My best guess is discarding
	// some of the tangential functions / interfaces would speed it up; but that in
	// of itself would be an interesting datapoint.

	b.Run("struct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := rangeStruct(0, 1024)
			var _ = r.toList()
		}
	})

	b.Run("closure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := rangeClosure(0, 1024)
			var _ = r.toList()
		}
	})
}

func BenchmarkRangeEach(b *testing.B) {

}
