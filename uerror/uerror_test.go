package uerror

import (
	"fmt"
	"io"
)

func ExampleWithError() {
	var err error
	err = WithError(io.ErrShortBuffer)
	err = WithErrorAndMessageF(err, "read %v end", "log.txt")
	err = WithErrorAndMessage(err, "b error")

	fmt.Println(err)
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
