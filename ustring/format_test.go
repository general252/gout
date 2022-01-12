package ustring

import (
	"fmt"
)

func ExampleFormat() {
	fmt.Println(Format())
	fmt.Println(Format(1))
	fmt.Println(Format(1, 2, 3))
	fmt.Println(Format("%d %v %+v", 1, 2, map[int]string{3: "hello"}))
	fmt.Println(Format("%#v", map[int]string{4: "world"}))

	// output:
	// 1
	// 1 2 3
	// 1 2 map[3:hello]
	// map[int]string{4:"world"}
}
