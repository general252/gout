package ucontainer

import "log"

func ExampleList_Delete() {
	l := NewList[int]()

	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(3)
	l.PushBack(3)
	l.PushBack(3)
	l.PushBack(4)

	l.Range(func(v int) bool {
		log.Println(v)
		return true
	})

	l.Delete(4, func(a, b int) bool {
		return a == b
	})

	log.Println("---------------------------")
	l.Range(func(v int) bool {
		log.Println(v)
		return true
	})

	// output:

}
