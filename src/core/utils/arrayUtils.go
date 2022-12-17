package utils

func Chunks[T any](array []T, chunkSize int) [][]T {
	steps := len(array) / chunkSize

	c := make([][]T, chunkSize)
	for i := 0; i < chunkSize; i++ {
		c[i] = array[(i * steps):((i + 1) * steps)]
	}
	return c
}
