package ucontainer

import "log"

func ExampleNewPool() {
	p := NewPoolWithSize(32)

	for i := 0; i < 66; i++ {
		if i%32 == 0 {
			log.Println()
		}

		if buf, err := p.Get(1); err != nil {
			log.Println(err)
		} else {
			log.Printf("%5d %p", i, buf)
		}
	}

	// output:
}
