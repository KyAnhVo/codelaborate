package main

type Ordered[T any] interface {
	// Compare compares this with other,
	// returns  1 iff this > other,
	// returns  0 iff this == other,
	// returns -1 iff this < other
	Compare(other T) int
}

const MinQueueCapacity int = 32

type Queue[T any] struct {
	// actually array
	slice 		[]T

	// dequeue from here
	head 		int

	// enqeue from here
	tail 		int
}

// NewQueue creates a queue of type T with the
// given capacity. If given capacity is less than MinQueueCapacity,
// we will use default MinQueueCapacity for capacity.
func NewQueue[T any](capacity int) *Queue[T] {
	if capacity < MinQueueCapacity {
		capacity = MinQueueCapacity
	}
	q := new(Queue[T])
	q.head 		= 0
	q.tail 		= 0
	q.slice 	= make([]T, capacity)
	return q
}



func (q *Queue[T]) Enqueue(v T) {
	if (q.tail + 1) % cap(q.slice) == q.head {
		newSliceInd := 0
		newSlice := make([]T, 2 * cap(q.slice))
		for q.head != q.tail {
			newSlice[newSliceInd] = q.slice[q.head]
			q.head = (q.head + 1) % cap(q.slice)
			newSliceInd++
		}
		q.head = 0
		q.tail = newSliceInd
		q.slice = newSlice
	}

	q.slice[q.tail] = v
	q.tail = (q.tail + 1) % cap(q.slice)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}
	val := q.slice[q.head]
	q.head = (q.head + 1) % cap(q.slice)

	// resize when go below 1/3 cap
	size := (q.tail + cap(q.slice) - q.head) % cap(q.slice)
	if cap(q.slice) >= MinQueueCapacity * 2 && size < cap(q.slice) / 3 {
		newSliceInd := 0
		newSlice := make([]T, cap(q.slice) / 2)
		for q.head != q.tail {
			newSlice[newSliceInd] = q.slice[q.head]
			q.head = (q.head + 1) % cap(q.slice)
			newSliceInd++
		}
		q.head = 0
		q.tail = newSliceInd
		q.slice = newSlice
	}

	return val, true
}

func (q *Queue[T]) IsEmpty() bool {
	return q.head == q.tail
}

