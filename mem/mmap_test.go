package mem

import "testing"

func TestAlloc(t *testing.T) {

	m, _ := SysAlloc(64)

	if len(*(*[]byte)(m)) != 64 {
		t.Errorf("Expected length 64, got %d", len(*(*[]byte)(m)))
	}

	SysFree(*(*[]byte)(m))

}
