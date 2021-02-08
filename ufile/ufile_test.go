package ufile

import "fmt"

func ExampleSplit() {
	dir, file, ext := Split("/home/name/hello.go")
	fmt.Println(dir, file, ext)
	// output:
	// /home/name/ hello .go
}
