package ef

type Atomic[V any] struct {
}

type AtomicNum[N Number] struct {
	v N
}
