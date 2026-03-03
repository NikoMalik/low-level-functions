//go:build amd64

package lowlevelfunctions

//go:noescape
func SumUnroll8(arr []int) int
