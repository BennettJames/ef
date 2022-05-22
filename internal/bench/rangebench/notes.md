
# Range Benchmarks


### Struct Vs Closure

This test compared defining a distinct struct to manage a range stream, vs
using a function closure.

A function closure is a little more convenient to use - the behavior is simple,
and the entire range can be defined within just a single function and ~10 lines.
The struct approach involves about three different declarations - still fairly
simple, but a little noiser than just doing it all inline.

Results (taken 2022-05-21):

```
goos: darwin
goarch: amd64
pkg: github.com/bennettjames/ef/internal/bench/rangebench
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkRangeStructVsClosure
BenchmarkRangeStructVsClosure/struct
BenchmarkRangeStructVsClosure/struct-16         	  153198	      7635 ns/op	   25232 B/op	      13 allocs/op
BenchmarkRangeStructVsClosure/closure
BenchmarkRangeStructVsClosure/closure-16        	  128311	      9053 ns/op	   25272 B/op	      15 allocs/op
PASS
ok  	github.com/bennettjames/ef/internal/bench/rangebench	2.705s
```

This was pretty consistent - the struct took ~84% of the time that the function
did.
