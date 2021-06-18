package ureader

import (
	"fmt"
	"io"
)

// processReader 获取进度辅助
type processReader struct {
	totalSize int64
	readSize  int64
	progress  Progress

	reader io.Reader
}

// Progress 进度
type Progress func(progress float64, readSize int64, totalSize int64)

// ProcessReader
func ProcessReader(r io.Reader, totalSize int64, progress Progress) io.Reader {
	return &processReader{
		totalSize: totalSize,
		readSize:  0,
		progress:  progress,
		reader:    r,
	}
}

func (c *processReader) Read(p []byte) (n int, err error) {
	if c.reader == nil {
		return -1, fmt.Errorf("from error")
	}

	n, err = c.reader.Read(p)
	if err == nil && c.progress != nil && c.totalSize > 0 {
		c.readSize += int64(n)
		var v = float64(c.readSize) / float64(c.totalSize) * float64(100)
		c.progress(v, c.readSize, c.totalSize)
	}

	return
}
