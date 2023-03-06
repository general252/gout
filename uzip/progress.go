package uzip

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Progress 进度
type Progress func(progress float64, readSize int64, totalSize int64)

func NewProcessReader(progress Progress) *processReader {
	return &processReader{
		totalSize: 0,
		readSize:  0,
		progress:  progress,
		reader:    nil,
	}
}

// processReader 获取进度辅助
type processReader struct {
	totalSize int64
	readSize  int64
	progress  Progress

	reader   io.Reader
	readerAt io.ReaderAt
}

func (c *processReader) SetReader(reader io.Reader) {
	c.reader = reader
}
func (c *processReader) SetReaderAt(readerAt io.ReaderAt) {
	c.readerAt = readerAt
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
func (c *processReader) ReadAt(p []byte, off int64) (n int, err error) {
	if c.readerAt == nil {
		return -1, fmt.Errorf("from error")
	}

	n, err = c.readerAt.ReadAt(p, off)
	if err == nil && c.progress != nil && c.totalSize > 0 {
		c.readSize += int64(n)
		var v = float64(c.readSize) / float64(c.totalSize) * float64(100)
		c.progress(v, c.readSize, c.totalSize)
	}

	return
}

// getSamplePart 获取数组中相同部分
func getSamplePart(arr []string) string {
	if len(arr) <= 0 {
		return ""
	}
	if len(arr) == 1 {
		return arr[0]
	}

	var minLength = 20480
	for _, s := range arr {
		if len(s) < minLength {
			minLength = len(s)
		}
	}

	var n = 0

	var end = false
	for i := 0; i < minLength; i++ {
		n = i

		ch := arr[0][i]
		for _, s := range arr {
			if ch != s[i] {
				end = true
				break
			}
		}

		if end {
			break
		}
	}

	return arr[0][:n]
}

func getBasePath(files []string) (string, error) {
	var commonPrefix string

	for i, file := range files {
		dir, err := filepath.Abs(filepath.Dir(file))
		if err != nil {
			return "", err
		}

		// Split the directory path into its components.
		parts := strings.Split(dir, string(filepath.Separator))

		// If this is the first file, just use its directory as the common prefix.
		if i == 0 {
			commonPrefix = dir
			continue
		}

		// Compare the directory path to the current common prefix.
		// Update the common prefix to include only the common parts.
		for j, part := range strings.Split(commonPrefix, string(filepath.Separator)) {
			if j >= len(parts) || part != parts[j] {
				commonPrefix = strings.Join(parts[:j], string(filepath.Separator))
				break
			}
		}
	}

	return commonPrefix, nil
}

// getTotalSize 获取文件大小
func getTotalSize(files []string) (int64, error) {
	var size int64
	for _, file := range files {
		err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				size += info.Size()
			}
			return err
		})
		if err != nil {
			return 0, err
		}
	}

	return size, nil
}

// createFile 创建文件或目录
func createFile(name string, isDir bool) (*os.File, error) {
	name = filepath.Clean(name)
	if isDir {
		err := os.MkdirAll(name, 0644)
		if err != nil {
			return nil, err
		}
		return nil, nil
	} else {
		err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, string(os.PathSeparator))]), 0644)
		if err != nil {
			return nil, err
		}
	}

	return os.Create(name)
}
