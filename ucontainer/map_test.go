package ucontainer

import "log"

func ExampleNewMap() {
	type A struct {
		ID   string
		Name string
	}

	m := NewMap[string, *A]()

	m.Store("a", &A{
		ID:   "1",
		Name: "tina",
	})
	m.Store("a", &A{
		ID:   "1",
		Name: "tina2",
	})
	m.Store("b", &A{
		ID:   "3",
		Name: "tony",
	})

	m.Range(func(k string, v *A) bool {
		log.Printf(">> k: %v, v: %v", k, v)
		return true
	})

	object, ok := m.Load("b")
	log.Println(object, ok)

	object, ok = m.Load("c")
	log.Println(object, ok)

	m.Delete("b")

	object, ok = m.Load("b")
	log.Println(object, ok)

	// output:
}
