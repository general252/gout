package ucontainer

import "log"

func ExampleNewArray() {
	type A struct {
		Name string
	}

	m := NewArray[*A]()
	m.Add(&A{Name: "1"})
	m.Add(&A{Name: "2"})
	m.Add(&A{Name: "3"})
	m.Delete(&A{Name: "2"}, func(a *A, b *A) bool {
		return a.Name == b.Name
	})

	m.Range(func(idx int, o *A) bool {
		log.Println(idx, o)
		return true
	})

	m.DeleteByIndex(0)

	m.Range(func(idx int, o *A) bool {
		log.Println(idx, o)
		return true
	})

	// output:
}
