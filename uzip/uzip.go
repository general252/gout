package uzip

import (
	"fmt"
	"strings"
)

// Compress 压缩 支持 *.tag.gz/*.zip
func Compress(files []string, dest string, progress Progress) error {
	if strings.HasSuffix(dest, ".tar.gz") {
		return GzipCompress(files, dest, progress)
	} else if strings.HasSuffix(dest, ".zip") {
		return ZipCompress(files, dest, progress)
	} else {
		return fmt.Errorf("not support")
	}
}

// DeCompress 解压 支持 *.tag.gz/*.zip
func DeCompress(file, dest string, progress Progress) error {
	if strings.HasSuffix(file, ".tar.gz") {
		return GzipDeCompress(file, dest, progress)
	} else if strings.HasSuffix(dest, ".zip") {
		return ZipDeCompress(file, dest, progress)
	} else {
		return fmt.Errorf("not support")
	}
}
