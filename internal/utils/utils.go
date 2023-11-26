package utils

// Map returns a slice of converted instances from T1 to T2 using the `mapper` function.`
func Map[T1, T2 any](items []T1, mapper func(item T1) T2) []T2 {
	result := make([]T2, 0, len(items))
	for _, item := range items {
		result = append(result, mapper(item))
	}

	return result
}
