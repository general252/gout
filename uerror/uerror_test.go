package uerror

import (
	"fmt"
	"log"
)

func acc() {
	log.Println()
	err1 := fmt.Errorf("hello error")

	err := PrintlnWithError(err1, 123, "124")
	_ = err
	//fmt.Println(err)
	var x, ok = ConvertToUError(err)
	fmt.Printf("-------------[%v][%v]--------------", x, ok)
}
func add() {
	acc()
}
func ExamplePrintfWithError() {
	add()
	// Output:
	//  error: hello error
	//  message: 123124
	//  stack:
	//   01. uerror_test.go:12
	//   02. uerror_test.go:19
	//   03. uerror_test.go:22
	//   04. run_example.go:62
	//   05. example.go:44
}
