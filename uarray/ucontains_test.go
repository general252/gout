// go test . -v  -count=1
// go test -bench=.
package uarray

import "testing"

func TestContains(t *testing.T) {
	var arr = []string{"aa", "bb", "cc"}
	var s = "bb"
	i := Contains(arr, s)
	if i != 1 {
		t.Error(i)
	}
}

func TestContainsFunction(t *testing.T) {
	var arr = []string{"aa", "bb", "cc"}
	var s = "bb"
	i := ContainsFunction(arr, s, func(a, b interface{}) bool {
		return a.(string) == b.(string)
	})
	if i != 1 {
		t.Error(i)
	}
}

func TestContainsString(t *testing.T) {
	var arr = []string{"aa", "bb", "cc"}
	var s = "bb"
	i := ContainsString(arr, s)
	if i != 1 {
		t.Error(i)
	}
}

func BenchmarkContains(b *testing.B) {
	sa := []string{"q", "w", "e", "r", "t"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Contains(sa, "r")
	}
}

func BenchmarkContainsString(b *testing.B) {
	sa := []string{"q", "w", "e", "r", "t"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ContainsString(sa, "r")
	}
}
