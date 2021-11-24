package generator

import "testing"

func TestDo(t *testing.T) {
	result1 := Do()
	result2 := Do()
	if len(result1) != shortLength || len(result2) != shortLength {
		t.Errorf("len(Do()) == %v, want %v", len(Do()), shortLength)
	}
	if result1 == result2 {
		t.Errorf("result1==result2, %v, %v", result1, result2)
	}
}
