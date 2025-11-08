package utils

func Contains[T comparable](src []T, dst T) bool {
	return ContainsFunc[T](src, func(src T) bool {
		return src == dst
	})
}

func ContainsFunc[T any](src []T, equal func(src T) bool) bool {
	// 遍历调用equal函数进行判断
	for _, v := range src {
		if equal(v) {
			return true
		}
	}
	return false
}

// Map slice to another
func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func AnyArrayToInterfaceArray[T any](items []T) []any {
	out := make([]any, len(items))
	for i, v := range items {
		out[i] = v
	}

	return out
}

// MoveToFront 将slice中的元素移动到首位
func MoveToFront[T any](slice []T, index int) []T {
	temp := slice[index]
	copy(slice[1:index+1], slice[0:index])
	slice[0] = temp
	return slice
}

// ArrayChunk 将slice分组
func ArrayChunk[T any](slice []T, size int) [][]T {
	chunks := make([][]T, 0)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// StringChunkWithHashPartation 将slice分组, 根据hash值分组
func StringChunkWithHashPartation(slice []string, size int) [][]string {
	if size <= 0 {
		return [][]string{}
	}

	chunks := make([][]string, size)
	for _, v := range slice {
		partition := HashToPartition(v, size)
		chunks[partition] = append(chunks[partition], v)
	}
	return chunks
}
