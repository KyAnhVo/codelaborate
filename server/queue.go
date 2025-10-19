package main

import (
	"sync"
)

type Queue[T any] struct {
	// true if empty, false otherwise
	isEmpty 	bool

	// max #elements of before resize
	capacity 	uint64

	// first element index
	first 		uint64

	// last element index + 1
	last 		uint64

	// actual implemented slice
	slice 		[]T

	// mutex to access Queue
	lock 		sync.Mutex
}

func NewQueue[T any](capacity uint64) *Queue[T] {
	if capacity == 0 {
		capacity = 32
	}

	q := new(Queue[T])
	q.isEmpty 	= true
	q.capacity 	= capacity
	q.first 	= 0
	q.last 		= 0
	q.slice 	= make([]T, q.capacity)

	return q
}

func (q *Queue[T]) Enqueue(v T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.first == q.last && !(q.isEmpty) { // queue is full
		newSlice 		:= make([]T, q.capacity * 2)
		var next uint64
		for next = 0; next < q.capacity; next++ {
			newSlice[next] = q.slice[q.first]
			q.first = (q.first + 1) % q.capacity
		}
		q.slice = newSlice
		q.first = 0
		q.last 	= next
		q.capacity *= 2
	}
	
	q.slice[q.last] = v
	q.last 			= (q.last + 1) % q.capacity
	q.isEmpty 		= false
}

func (q *Queue[T]) Dequeue() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.isEmpty {
		var zero T
		return zero, false
	}

	val := q.slice[q.first]
	q.first = (q.first + 1) % q.capacity
	q.isEmpty = q.first == q.last
	return val, true
}

func (q *Queue[T]) IsEmpty() bool {
	return q.isEmpty
}
