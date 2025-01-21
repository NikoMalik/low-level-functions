package lowlevelfunctions

import (
	"os"
	"sync"
)

var mutableStringPool = &sync.Pool{
	New: func() any {
		buf := MakeNoZeroCap(0, int(StrSize*uintptr(os.Getpagesize())))
		return (*MutableString)(&buf)
	},
}

// for big changed strings
func AcquireMutableString() *MutableString {
	return mutableStringPool.Get().(*MutableString)
}

func ReleaseMutableString(m *MutableString) {
	*m = (*m)[:0]
	mutableStringPool.Put(m)
}
