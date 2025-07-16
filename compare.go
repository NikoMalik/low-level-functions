package lowlevelfunctions

import "unsafe"

//go:noescape
func CompareImpl(a, b unsafe.Pointer, len int) bool
