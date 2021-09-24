package ushell

import (
	"bytes"
	"context"
	"github.com/axgle/mahonia"
	"io"
	"log"
	"os/exec"
	"strings"
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

	_ = ShellCommandStreamV2(context.TODO(), "ls -l", func(stdoutPipe io.ReadCloser, stderrPipe io.ReadCloser) {
	}, func(c *exec.Cmd) {
		c.Stdin = buf
	})

	// output:
	//
}

func ExampleShellCommandStreamArgs() {
	mysqlDump := `C:\Program Files (x86)\mysql\bin\mysqldump.exe`
	args := strings.Split("-v --host=127.0.0.1 --port=3306 --user=root --password=123456 --compress --databases mysql", " ")
	args = append(args, ">")
	args = append(args, `C:\a b\backup.sql`)

	var dec = mahonia.NewDecoder("gbk")
	err := ShellCommandStreamArgs(context.TODO(), mysqlDump, args, func(msg []byte, isError bool) {
		var m = dec.ConvertString(string(msg))
		log.Println(m)
	})
	log.Println(err)

	// output:
	//
}
