package set

// Diff returns a set containing every element from two hash sets that does not exist in both
func Diff[T comparable](first, second *HashSet[T]) *HashSet[T] {
	first.Lock()
	defer first.Unlock()
	second.Lock()
	defer second.Unlock()

	diff := make(map[T]struct{})

	for value, _ := range first.underlying {
		if _, ok := second.underlying[value]; !ok {
			diff[value] = struct{}{}
		}
	}

	for value, _ := range second.underlying {
		if _, ok := first.underlying[value]; !ok {
			diff[value] = struct{}{}
		}
	}

	return &HashSet[T]{
		underlying: diff,
	}
}
