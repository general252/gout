package uzip

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Compress 压缩 使用gzip压缩成tar.gz
func Compress(files []string, dest string, progress Progress) error {
	if files == nil || len(files) == 0 {
		return fmt.Errorf("nof file")
	}

	d, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer d.Close()

	gw := gzip.NewWriter(d)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	var newFiles []string
	for _, file := range files {
		newFiles = append(newFiles, filepath.Clean(file))
	}

	progressReader := NewProcessReader(progress)
	if progressReader.totalSize, err = getTotalSize(newFiles); err != nil {
		return err
	}

	basePath := getSamplePart(newFiles)
	if len(newFiles) == 1 {
		basePath = filepath.Dir(basePath)
		basePath = filepath.Clean(basePath) + string(os.PathSeparator)
	} else {
		basePath = filepath.Clean(basePath)
	}

	for _, file := range newFiles {
		err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return err
			}

			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}

			header.Name = strings.TrimPrefix(filepath.Clean(path), basePath)
			if err = tw.WriteHeader(header); err != nil {
				return err
			}

			if f, err := os.Open(path); err != nil {
				return err
			} else {
				defer f.Close()

				progressReader.SetReader(f)

				if _, err = io.Copy(tw, progressReader); err != nil {
					return err
				}
			}

			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// DeCompress 解压 tar.gz
func DeCompress(tarFile, dest string, progress Progress) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	progressReader := NewProcessReader(progress)
	if info, err := srcFile.Stat(); err != nil {
		return err
	} else {
		progressReader.totalSize = info.Size()
		progressReader.SetReader(srcFile)
	}

	gr, err := gzip.NewReader(progressReader)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		filename := fmt.Sprintf("%v/%v", dest, hdr.Name)

		file, err := createFile(filename)
		if err != nil {
			return err
		}

		_, _ = io.Copy(file, tr)
		file.Close()
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	name = filepath.Clean(name)
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, string(os.PathSeparator))]), 0644)
	if err != nil {
		return nil, err
	}

	return os.Create(name)
}

// processReader 获取进度辅助
type processReader struct {
	totalSize int64
	readSize  int64
	progress  Progress

	reader io.Reader
}

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

func (c *processReader) SetReader(reader io.Reader) {
	c.reader = reader
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

// getSamplePart 获取数组中相同部分
func getSamplePart(arr []string) string {
	if len(arr) <= 0 {
		return ""
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
