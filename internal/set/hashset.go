package set

import "sync"

// HashSet represents a simple thread safe hash set
type HashSet[T comparable] struct {
	sync.RWMutex
	underlying map[T]struct{}
}

// NewHashSet creates a new thread safe hash set
func NewHashSet[T comparable]() *HashSet[T] {
	return &HashSet[T]{
		underlying: make(map[T]struct{}),
	}
}

// Size returns the amount of elements stored inside the set
func (set *HashSet[T]) Size() int {
	return len(set.underlying)
}

// Contains checks if a specific value is present in the set
func (set *HashSet[T]) Contains(value T) bool {
	_, ok := set.underlying[value]
	return ok
}

// Add adds a value to the set
func (set *HashSet[T]) Add(value T) {
	set.Lock()
	defer set.Unlock()
	set.underlying[value] = struct{}{}
}

// ToSlice returns a slice containing every element of the set
func (set *HashSet[T]) ToSlice() []T {
	set.Lock()
	defer set.Unlock()

	values := make([]T, 0, len(set.underlying))
	for value, _ := range set.underlying {
		values = append(values, value)
	}
	return values
}
