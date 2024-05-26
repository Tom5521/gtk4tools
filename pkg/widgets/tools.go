package widgets

func splice[T any](start, count int, items []T) (ret []T) {
	ret = make([]T, len(items)-count)
	copy(ret, items[:start])
	copy(ret[start:], items[start+count:])
	return
}
