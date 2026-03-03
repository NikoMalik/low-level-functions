package lowlevelfunctions

import (
	_ "unsafe" // for go:linkname
)

// Prefetch prefetches data from memory to cache using compiler intrinsic.
// This is linked to internal/runtime/sys.Prefetch which compiles to inline
// PREFETCH instruction without function call overhead.
//
// AMD64: Produces PREFETCHT0 instruction (inline)
// ARM64: Produces PRFM instruction with PLDL1KEEP option (inline)
//
//go:linkname prefetch internal/runtime/sys.Prefetch
func prefetch(addr uintptr)
