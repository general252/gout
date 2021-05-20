package ushell

import (
	"bytes"
	"context"
	"io"
	"log"
	"os/exec"
)

type myBuffer struct {
	buf *bytes.Buffer
}

func (c *myBuffer) Read(p []byte) (n int, err error) {
	n, err = c.buf.Read(p)
	log.Println(n)
	return n, err
}

func ExampleShellCommandStreamV2() {
	var buf = &myBuffer{buf: bytes.NewBufferString("hello")}

	_ = ShellCommandStreamV2(context.TODO(), "ls -l", func(c *exec.Cmd, stdoutPipe io.ReadCloser, stderrPipe io.ReadCloser) {
		c.Stdin = buf
	})

	// output:
	//
}
