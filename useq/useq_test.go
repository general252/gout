package useq

import (
	"fmt"
	"sync"
	"testing"
)

func ExampleGetSeqU64() {
	var seq = GetSeqU64()
	fmt.Println(seq)

	seq = GetSeqU64()
	fmt.Println(seq)

	// Output:
	// 1
	// 2
}

func ExampleGetSeq32() {
	var wg = sync.WaitGroup{}
	var fun = func(id int) {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			var seq = GetSeq32()
			fmt.Printf("%4d %5d\n", id, seq)
		}
	}

	var count = 4
	wg.Add(count)
	for i := 0; i < count; i++ {
		go fun(i + 1)
	}
	wg.Wait()
	// Output:
	// 1
}

func BenchmarkGetSeq32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var seq = GetSeq64()
		_ = seq
	}
}
