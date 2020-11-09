package uerror

import (
	"fmt"
	"io"
)

func ExamplePrintfWithError() {
	var a = func() error {
		return WithErrorAndMessageF(io.EOF, "read %v end", "log.txt")
	}
	var b = func() error {
		return WithErrorAndMessage(a(), "b error")
	}
	var c = func() error {
		return b()
	}
	fmt.Println(c())
	// Output:
	// error: EOF
	// message:
	//	 01. read log.txt end
	//   02. b error
	// stack:
	//   01. main.go:11
	//   02. main.go:14
	//   03. main.go:17
	//   04. main.go:19
	//   05. proc.go:203
}
