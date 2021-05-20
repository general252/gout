package uzip

import (
	"log"
)

func ExampleGetSamplePart() {
	x := getSamplePart([]string{
		"E:/a.txtx",
		"E:/a.txty",
		"E:/a.txtx",
	})
	log.Println(x)

	// output:
	//
}

func ExampleCompress() {
	_ = Compress([]string{"E://backup", "E:/c.txt", "E:\\\\test\\main.go"}, "E:/abc.gzip", func(progress float64, readSize int64, totalSize int64) {
		log.Println(progress, readSize, totalSize)
	})

	// output:
	//
}

func ExampleDeCompress() {
	err := DeCompress("E:/abc.gzip", "E:/tmp", func(progress float64, readSize int64, totalSize int64) {
		log.Println(progress, readSize, totalSize)
	})
	log.Println(err)

	// output:
	//
}
