package queue

import (
	"container/list"
	"sync"
)

// Queue represents a simple thread safe FIFO queue
type Queue[T any] struct {
	sync.RWMutex
	queue *list.List
}

// New creates a new queue
func New[T any]() *Queue[T] {
	return &Queue[T]{
		queue: list.New(),
	}
}

// Size returns the amount of queued entries
func (queue *Queue[T]) Size() int {
	queue.Lock()
	defer queue.Unlock()
	return queue.queue.Len()
}

// Push pushes entries to the back of the queue
func (queue *Queue[T]) Push(values ...T) {
	queue.Lock()
	defer queue.Unlock()
	for _, val := range values {
		queue.queue.PushBack(val)
	}
}

func (queue *Queue[T]) unsafePop() (T, bool) {
	if queue.queue.Len() == 0 {
		var zero T
		return zero, false
	}
	elem := queue.queue.Front()
	queue.queue.Remove(elem)
	return elem.Value.(T), true
}

// Pop pops the first element of the queue
func (queue *Queue[T]) Pop() (T, bool) {
	queue.Lock()
	defer queue.Unlock()
	return queue.unsafePop()
}

// PopN pops the first n elements of the queue
func (queue *Queue[T]) PopN(n int) []T {
	queue.Lock()
	defer queue.Unlock()
	values := make([]T, 0, n)
	for i := 0; i < n; i++ {
		val, ok := queue.unsafePop()
		if !ok {
			break
		}
		values = append(values, val)
	}
	return values
}
