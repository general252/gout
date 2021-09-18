package uzip

import (
	"fmt"
	"strings"
)

// Compress 压缩 使用gzip压缩成tar.gz
func Compress(files []string, dest string, progress Progress) error {
	if strings.HasSuffix(dest, ".tar.gz") {
		return GzipCompress(files, dest, progress)
	} else if strings.HasSuffix(dest, ".zip") {
		return ZipCompress(files, dest, progress)
	} else {
		return fmt.Errorf("not support")
	}
}

// DeCompress 解压 tar.gz
func DeCompress(file, dest string, progress Progress) error {
	if strings.HasSuffix(file, ".tar.gz") {
		return GzipDeCompress(file, dest, progress)
	} else {
		return fmt.Errorf("not support")
	}
}
