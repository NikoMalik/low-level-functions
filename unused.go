package lowlevelfunctions

func new[T any](value T, size int) *T {
	obj := (*T)(mallocgc(uintptr(size), nil, false))
	*obj = value
	return obj
}
