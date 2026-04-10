package lowlevelfunctions

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"unsafe"
)

// Node for linked list - realistic use case for prefetch
type Node struct {
	Value int64
	Next  *Node
	_     [56]byte // padding to fill cache line (64 bytes total)
}

// createLinkedList creates a linked list with random order
func createLinkedList(size int, randomize bool) *Node {
	if size == 0 {
		return nil
	}

	nodes := make([]Node, size)

	// Initialize nodes
	for i := 0; i < size; i++ {
		nodes[i].Value = int64(i)
	}

	// Create indices for linking
	indices := make([]int, size)
	for i := range indices {
		indices[i] = i
	}

	// Shuffle if randomize is true
	if randomize {
		rand.Shuffle(len(indices), func(i, j int) {
			indices[i], indices[j] = indices[j], indices[i]
		})
	}

	// Link nodes
	for i := 0; i < size-1; i++ {
		nodes[indices[i]].Next = &nodes[indices[i+1]]
	}
	nodes[indices[size-1]].Next = nil

	return &nodes[indices[0]]
}

// BenchmarkPrefetchLinkedList - realistic scenario where prefetch helps
func BenchmarkPrefetchLinkedList(b *testing.B) {
	sizes := []int{
		1000,    // ~64KB
		10000,   // ~640KB
		100000,  // ~6.4MB
		1000000, // ~64MB
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Sequential_%d/NoPrefetch", size), func(b *testing.B) {
			head := createLinkedList(size, false)
			runtime.GC() // ensure clean state
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				sum := int64(0)
				node := head
				for node != nil {
					sum += node.Value
					node = node.Next
				}
				runtime.KeepAlive(sum)
			}
		})

		b.Run(fmt.Sprintf("Sequential_%d/WithPrefetch", size), func(b *testing.B) {
			head := createLinkedList(size, false)
			runtime.GC()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				sum := int64(0)
				node := head
				for node != nil {
					// Prefetch next node while processing current
					if node.Next != nil {
						T0(unsafe.Pointer(node.Next))
					}
					sum += node.Value
					node = node.Next
				}
				runtime.KeepAlive(sum)
			}
		})

		b.Run(fmt.Sprintf("Random_%d/NoPrefetch", size), func(b *testing.B) {
			head := createLinkedList(size, true)
			runtime.GC()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				sum := int64(0)
				node := head
				for node != nil {
					sum += node.Value
					node = node.Next
				}
				runtime.KeepAlive(sum)
			}
		})

		b.Run(fmt.Sprintf("Random_%d/WithPrefetch", size), func(b *testing.B) {
			head := createLinkedList(size, true)
			runtime.GC()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				sum := int64(0)
				node := head
				for node != nil {
					if node.Next != nil {
						T0(unsafe.Pointer(node.Next))
					}
					sum += node.Value
					node = node.Next
				}
				runtime.KeepAlive(sum)
			}
		})

		b.Run(fmt.Sprintf("Random_%d/WithDoublePrefetch", size), func(b *testing.B) {
			head := createLinkedList(size, true)
			runtime.GC()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				sum := int64(0)
				node := head
				for node != nil {
					// Prefetch 2 nodes ahead
					if node.Next != nil {
						T0(unsafe.Pointer(node.Next))
						if node.Next.Next != nil {
							T0(unsafe.Pointer(node.Next.Next))
						}
					}
					sum += node.Value
					node = node.Next
				}
				runtime.KeepAlive(sum)
			}
		})
	}
}

func BenchmarkPrefetchArrayTraversal(b *testing.B) {
	const size = 8 * 1024 * 1024 // 8MB

	indices := make([]int, size/64)
	for i := range indices {
		indices[i] = i * 64
	}
	rand.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	b.Run("RandomAccess/NoPrefetch", func(b *testing.B) {
		data := make([]byte, size)
		// Touch data to make it cold
		for i := range data {
			data[i] = byte(i)
		}
		runtime.GC()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sum := byte(0)
			for _, idx := range indices {
				sum += data[idx]
			}
			runtime.KeepAlive(sum)
		}
	})

	b.Run("RandomAccess/WithPrefetch_4ahead", func(b *testing.B) {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i)
		}
		runtime.GC()

		const prefetchAhead = 4
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sum := byte(0)
			for j, idx := range indices {
				// Prefetch 4 accesses ahead
				if j+prefetchAhead < len(indices) {
					T0(unsafe.Pointer(&data[indices[j+prefetchAhead]]))
				}
				sum += data[idx]
			}
			runtime.KeepAlive(sum)
		}
	})

	b.Run("RandomAccess/WithPrefetch_8ahead", func(b *testing.B) {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i)
		}
		runtime.GC()

		const prefetchAhead = 8
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sum := byte(0)
			for j, idx := range indices {
				if j+prefetchAhead < len(indices) {
					T0(unsafe.Pointer(&data[indices[j+prefetchAhead]]))
				}
				sum += data[idx]
			}
			runtime.KeepAlive(sum)
		}
	})

	b.Run("RandomAccess/WithPrefetch_16ahead", func(b *testing.B) {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i)
		}
		runtime.GC()

		const prefetchAhead = 16
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sum := byte(0)
			for j, idx := range indices {
				if j+prefetchAhead < len(indices) {
					T0(unsafe.Pointer(&data[indices[j+prefetchAhead]]))
				}
				sum += data[idx]
			}
			runtime.KeepAlive(sum)
		}
	})
}

// BenchmarkPrefetchMatrixTraversal - 2D array access patterns
func BenchmarkPrefetchMatrixTraversal(b *testing.B) {
	const size = 2048
	matrix := make([][]int32, size)
	for i := range matrix {
		matrix[i] = make([]int32, size)
		for j := range matrix[i] {
			matrix[i][j] = int32(i*size + j)
		}
	}

	// Column-major access (cache-unfriendly)
	b.Run("ColumnMajor/NoPrefetch", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			sum := int32(0)
			for col := 0; col < size; col++ {
				for row := 0; row < size; row++ {
					sum += matrix[row][col]
				}
			}
			runtime.KeepAlive(sum)
		}
	})

	b.Run("ColumnMajor/WithPrefetch", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			sum := int32(0)
			for col := 0; col < size; col++ {
				for row := 0; row < size; row++ {
					// Prefetch next row in same column
					if row+8 < size {
						T0(unsafe.Pointer(&matrix[row+8][col]))
					}
					sum += matrix[row][col]
				}
			}
			runtime.KeepAlive(sum)
		}
	})
}
