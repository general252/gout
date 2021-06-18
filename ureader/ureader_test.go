package ureader

import (
	"bytes"
	"io"
	"log"
)

func ExampleNewProcessReader() {
	limitSize := int64(64)
	r := bytes.NewReader(bytes.Repeat([]byte("a"), int(limitSize)))

	progress := ProcessReader(r, limitSize, func(progress float64, readSize int64, totalSize int64) {
		log.Println(progress, readSize, totalSize)
	})

	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, progress)
	log.Println(err, buffer.String())

	// output:
	//
}
