package lowlevelfunctions

import (
	"os"
	"sync"
)

var mutableStringPool = &sync.Pool{
	New: func() any {
		buf := MakeNoZeroCap(0, StrSize*os.Getpagesize())
		return (*MutableString)(&buf)
	},
}

func AcquireMutableString() *MutableString {
	return mutableStringPool.Get().(*MutableString)
}

func ReleaseMutableString(m *MutableString) {
	*m = (*m)[:0]
	mutableStringPool.Put(m)
}
