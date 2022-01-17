package ucontainer

import (
	"log"
)

func ExampleNewQueue() {
	q := NewQueue()

	for i := 0; i < 10; i++ {
		q.Push(i)
	}

	for {
		v, ok := q.Pop()
		if !ok {
			break
		}

		log.Println(q.Len(), v)
	}

	// output:
}
