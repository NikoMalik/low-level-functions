//go:build amd64

package lowlevelfunctions

import "unsafe"

//go:noescape
func CompareImpl(a, b unsafe.Pointer, len int) bool

//go:noescape
func containsByteSliceAVX2(
	slicePtr unsafe.Pointer,
	sliceLen int,
	valuePtr unsafe.Pointer,
	valueLen int,
) bool

//go:noecape
func containsStringAVX2(
	slicePtr unsafe.Pointer,
	sliceLen int,
	valuePtr unsafe.Pointer,
	valueLen int,
) bool
