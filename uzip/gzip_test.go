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
	_ = Compress([]string{"E:\\test\\c.txt", "E:\\test\\backup"}, "E:\\test\\test.tar.gz", func(progress float64, readSize int64, totalSize int64) {
		log.Println(progress, readSize, totalSize)
	})

	// output:
	//
}

func ExampleCompress2() {
	_ = Compress([]string{"E:\\test\\backup\\a.txt"}, "E:\\test\\a.tar.gz", func(progress float64, readSize int64, totalSize int64) {
		log.Println(progress, readSize, totalSize)
	})

	// output:
	//
}

func ExampleDeCompress() {
	err := DeCompress("E:\\test\\test.tar.gz", "E:\\test\\out", func(progress float64, readSize int64, totalSize int64) {
		log.Println(progress, readSize, totalSize)
	})
	log.Println(err)

	// output:
	//
}

func ExampleCompress3() {
	err := DeCompress(`C:\Users\tony\Desktop\template.tar.gz`, `C:\Users\tony\Desktop\out`, func(progress float64, readSize int64, totalSize int64) {
		log.Println(progress, readSize, totalSize)
	})
	log.Println(err)

	// output:
	//
}