package helpers

func SliceFilter[T any](slice []T, test func(T) bool) []T {
	ret := make([]T, 0, len(slice))
	for _, s := range slice {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return ret
}
