package ef

import "sync/atomic"

type (
	// (todo)
	AtomicRef[V any] struct {
		v atomic.Value
	}

	// (todo)
	AtomicInt[N Integer] interface {
		// (todo)
		Get() N

		// (todo)
		Set(val N)

		// (todo)
		Add(val N) (new N)

		// (todo)
		CompareAndSwap(old, new N) (swapped bool)

		// (todo)
		Swap(new N) (old N)
	}

	atomicInt32 struct {
		val int32
	}

	atomicInt64 struct {
		val int64
	}

	atomicUint32 struct {
		val uint32
	}

	atomicUint64 struct {
		val uint64
	}

	atomicUintptr struct {
		val uintptr
	}
)

// (todo)
func (ar AtomicRef[V]) CompareAndSwap(old, new V) (swapped bool) {
	return ar.v.CompareAndSwap(old, new)
}

// (todo)
func (ar AtomicRef[V]) Load() (val V) {
	return ar.v.Load().(V)
}

// (todo)
func (ar AtomicRef[V]) Store(val V) {
	ar.v.Store(val)
}

// (todo)
func (ar AtomicRef[V]) Swap(new V) (old V) {
	return ar.v.Swap(new).(V)
}

// (todo)
func NewAtomicInt32(value int32) AtomicInt[int32] {
	return &atomicInt32{value}
}

func (ai *atomicInt32) Get() int32 {
	return atomic.LoadInt32(&ai.val)
}

func (ai *atomicInt32) Set(val int32) {
	atomic.StoreInt32(&ai.val, val)
}

func (ai *atomicInt32) Add(val int32) int32 {
	return atomic.AddInt32(&ai.val, val)
}

func (ai *atomicInt32) CompareAndSwap(old, new int32) (swapped bool) {
	return atomic.CompareAndSwapInt32(&ai.val, old, new)
}

func (ai *atomicInt32) Swap(new int32) (old int32) {
	return atomic.SwapInt32(&ai.val, new)
}

// (todo)
func NewAtomicInt64(value int64) AtomicInt[int64] {
	return &atomicInt64{value}
}

func (ai *atomicInt64) Get() int64 {
	return atomic.LoadInt64(&ai.val)
}

func (ai *atomicInt64) Set(val int64) {
	atomic.StoreInt64(&ai.val, val)
}

func (ai *atomicInt64) Add(val int64) int64 {
	return atomic.AddInt64(&ai.val, val)
}

func (ai *atomicInt64) CompareAndSwap(old, new int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(&ai.val, old, new)
}

func (ai *atomicInt64) Swap(new int64) (old int64) {
	return atomic.SwapInt64(&ai.val, new)
}

// (todo)
func NewAtomicUint32(value uint32) AtomicInt[uint32] {
	return &atomicUint32{value}
}

func (ai *atomicUint32) Get() uint32 {
	return atomic.LoadUint32(&ai.val)
}

func (ai *atomicUint32) Set(val uint32) {
	atomic.StoreUint32(&ai.val, val)
}

func (ai *atomicUint32) Add(val uint32) uint32 {
	return atomic.AddUint32(&ai.val, val)
}

func (ai *atomicUint32) CompareAndSwap(old, new uint32) (swapped bool) {
	return atomic.CompareAndSwapUint32(&ai.val, old, new)
}

func (ai *atomicUint32) Swap(new uint32) (old uint32) {
	return atomic.SwapUint32(&ai.val, new)
}

// (todo)
func NewAtomicUint64(value uint64) AtomicInt[uint64] {
	return &atomicUint64{value}
}

func (ai *atomicUint64) Get() uint64 {
	return atomic.LoadUint64(&ai.val)
}

func (ai *atomicUint64) Set(val uint64) {
	atomic.StoreUint64(&ai.val, val)
}

func (ai *atomicUint64) Add(val uint64) uint64 {
	return atomic.AddUint64(&ai.val, val)
}

func (ai *atomicUint64) CompareAndSwap(old, new uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64(&ai.val, old, new)
}

func (ai *atomicUint64) Swap(new uint64) (old uint64) {
	return atomic.SwapUint64(&ai.val, new)
}

// (todo)
func NewAtomicUintptr(value uintptr) AtomicInt[uintptr] {
	return &atomicUintptr{value}
}

func (ai *atomicUintptr) Get() uintptr {
	return atomic.LoadUintptr(&ai.val)
}

func (ai *atomicUintptr) Set(val uintptr) {
	atomic.StoreUintptr(&ai.val, val)
}

func (ai *atomicUintptr) Add(val uintptr) uintptr {
	return atomic.AddUintptr(&ai.val, val)
}

func (ai *atomicUintptr) CompareAndSwap(old, new uintptr) (swapped bool) {
	return atomic.CompareAndSwapUintptr(&ai.val, old, new)
}

func (ai *atomicUintptr) Swap(new uintptr) (old uintptr) {
	return atomic.SwapUintptr(&ai.val, new)
}
