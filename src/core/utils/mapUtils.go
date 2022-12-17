package utils

func MapCopy[K comparable, V any](src map[K]V, dest map[K]V) map[K]V {
	for k, v := range src {
		dest[k] = v
	}
	return dest
}
