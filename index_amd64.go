package lowlevelfunctions

//go:noescape
func indexAvx2(haystack []byte, needle []byte) int64
