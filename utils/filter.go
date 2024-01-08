package utils

func Filter[T any](slice []T, f func(T) bool) []T {
	var filtered []T

	for _, element := range slice {
		if f(element) {
			filtered = append(filtered, element)
		}
	}

	return filtered
}
