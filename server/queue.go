package main

type Ordered[T any] interface {
	// Compare compares this with other,
	// returns  1 iff this > other,
	// returns  0 iff this == other,
	// returns -1 iff this < other
	Compare(other T) int
}

const MinQueueCapacity uint64 = 32

type Queue[T Ordered[T]] struct {
	// true if empty, false otherwise
	isEmpty 	bool

	// max #elements of before resize
	capacity 	uint64

	// current size
	size		uint64

	// actual implemented slice
	slice 		[]T
}

// NewQueue creates a queue of type T with the
// given capacity. If given capacity is equal to 0,
// we will use default 32 capacity.
func NewQueue[T Ordered[T]](capacity uint64) *Queue[T] {
	if capacity < MinQueueCapacity {
		capacity = MinQueueCapacity
	}

	q := new(Queue[T])
	q.isEmpty 	= true
	q.capacity 	= capacity
	q.size 		= 0
	q.slice 	= make([]T, q.capacity)

	return q
}



func (q *Queue[T]) Enqueue(v T) {
	if q.size == q.capacity {
		newSlice := make([]T, 2 * q.capacity)
		for i := uint64(0); i < q.size; i++ {
			newSlice[i] = q.slice[i]
		}
		q.slice = newSlice
		q.capacity *= 2
	}

	// insert v at the end, then ensure heap property
	q.slice[q.size] = v
	ind := q.size
	for ind > 0 && q.slice[ind].Compare(q.slice[(ind - 1) / 2]) == -1 {
		swap(q.slice, ind, (ind - 1) / 2);
		ind = (ind - 1) / 2
	}

	// update metadata
	q.size += 1
	q.isEmpty = false
}

func (q *Queue[T]) Dequeue() (T, bool) {
	// nothing to dequeue
	if q.IsEmpty() {
		var zero T
		return zero, false
	}
	
	// save value for return
	retVal := q.slice[0]

	// swap last element with root
	swap(q.slice, 0, q.size - 1)
	q.size--

	// preserve heap invariant
	curr := uint64(0)
	for {
		left := 2 * curr + 1
		right := left + 1
		
		// no child
		if left >= q.size {
			break;
		}

		if right >= q.size { // one child
			if q.slice[left].Compare(q.slice[curr]) == -1 {
				swap(q.slice, curr, left)
				curr = left
				continue
			}

			break;
		} else { // two children
			min := left
			if q.slice[left].Compare(q.slice[right]) == 1 {
				min = right
			}
			if q.slice[min].Compare(q.slice[curr]) == -1 {
				swap(q.slice, curr, min)
				curr = min
				continue
			}

			break;
		}

		// all children not 
	}



	// metadata update
	q.isEmpty = q.size == 0

	// also updates capacity if goes small enough
	if q.capacity >= MinQueueCapacity * 2 && q.size < q.capacity / 3 {
		newSlice := make([]T, q.capacity / 2)
		for i := uint64(0); i < q.size; i++ {
			newSlice[i] = q.slice[i]
		}
		q.slice = newSlice
		q.capacity /= 2
	}

	return retVal, true
}

func (q *Queue[T]) IsEmpty() bool {
	return q.isEmpty
}

func swap[T any](arr []T, i1 uint64, i2 uint64) {
	if i1 == i2 {
		return
	}
	temp := arr[i1]
	arr[i1] = arr[i2]
	arr[i2] = temp
}
