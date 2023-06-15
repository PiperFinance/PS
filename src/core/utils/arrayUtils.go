package utils

func Chunks[T any](array []T, chunkSize int) [][]T {
	//var _size int
	//var steps int
	//if chunkSize > len(array) {
	//	_size = len(array)
	//	steps = 1
	//} else {
	//	_size = chunkSize
	//	steps = len(array) / chunkSize
	//}
	//c := make([][]T, _size)
	//for i := 0; i < _size; i++ {
	//	c[i] = array[(i * steps):((i + 1) * steps)]
	//}

	steps := (len(array) / chunkSize) + 1

	c := make([][]T, steps)
	for i := 0; i < steps; i++ {
		outerBound := ((i + 1) * chunkSize)
		if outerBound > len(array) {
			outerBound = len(array)
		}
		c[i] = array[(i * chunkSize):outerBound]
	}

	return c
}
