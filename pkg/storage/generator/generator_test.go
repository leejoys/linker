package generator

import "testing"

func TestDo(t *testing.T) {
	g := New(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_",
		10)
	result1 := g.Do()
	result2 := g.Do()
	if len(result1) != g.length || len(result2) != g.length {
		t.Errorf("len(Do()) == %v, want %v", len(g.Do()), g.length)
	}
	if result1 == result2 {
		t.Errorf("result1==result2, %v, %v", result1, result2)
	}
}
