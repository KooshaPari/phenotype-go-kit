package ringbuffer

// RingBuffer is a generic circular buffer that stores items of any type.
type RingBuffer[T any] struct {
	items     []T
	nextIndex int
	count     int
	size      int
}

// New creates a new ring buffer with the specified capacity.
func New[T any](size int) *RingBuffer[T] {
	return &RingBuffer[T]{
		items: make([]T, size),
		size:  size,
	}
}

// Push adds an item to the ring buffer, overwriting the oldest if full.
func (r *RingBuffer[T]) Push(item T) {
	r.items[r.nextIndex] = item
	r.nextIndex = (r.nextIndex + 1) % r.size
	if r.count < r.size {
		r.count++
	}
}

// GetAll returns all items in the buffer, oldest first.
func (r *RingBuffer[T]) GetAll() []T {
	result := make([]T, r.count)
	for i := 0; i < r.count; i++ {
		result[i] = r.items[(r.nextIndex-r.count+i+r.size)%r.size]
	}
	return result
}

// Len returns the number of items currently in the buffer.
func (r *RingBuffer[T]) Len() int {
	return r.count
}

// Cap returns the capacity of the ring buffer.
func (r *RingBuffer[T]) Cap() int {
	return r.size
}
