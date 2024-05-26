package widgets

func splice[T any](slice []T, start, count int, items []T) []T {
	newLength := len(slice) - count + len(items)
	result := make([]T, newLength)
	copy(result[:start], slice[:start])
	copy(result[start:start+len(items)], items)
	copy(result[start+len(items):], slice[start+count:])
	return result
}
