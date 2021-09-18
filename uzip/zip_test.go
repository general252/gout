package uzip

import "log"

func ExampleZipCompress() {
	err := ZipCompress([]string{"E:\\test\\backup\\a.txt", "E:/xxx/"}, "E:\\test\\a.zip", func(progress float64, readSize int64, totalSize int64) {
		//log.Println(progress, readSize, totalSize)
	})
	log.Println(err)

	// output:
	//
}

func ExampleZipCompress1() {
	err := ZipCompress([]string{"E:\\test\\backup\\a.txt"}, "E:\\test\\a.zip", func(progress float64, readSize int64, totalSize int64) {
		log.Println(progress, readSize, totalSize)
	})
	log.Println(err)

	// output:
	//
}

func ExampleZipDeCompress() {
	err := ZipDeCompress("E:\\test\\a.zip", "E:\\test\\out", func(progress float64, readSize int64, totalSize int64) {
		//log.Println(progress, readSize, totalSize)
	})
	log.Println(err)

	// output:
	//
}
