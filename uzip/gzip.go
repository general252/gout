package uzip

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/axgle/mahonia"
)

// GzipCompress 压缩 使用gzip压缩成tar.gz
func GzipCompress(files []string, dest string, progress Progress) error {
	if files == nil || len(files) == 0 {
		return fmt.Errorf("nof file")
	}

	d, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer d.Close()

	var newFiles []string
	for _, file := range files {
		newFiles = append(newFiles, filepath.Clean(file))
	}

	// 进度
	progressReader := NewProcessReader(progress)
	if progressReader.totalSize, err = getTotalSize(newFiles); err != nil {
		return err
	}

	// 计算顶层目录
	basePath := getSamplePart(newFiles)
	{
		if len(newFiles) == 1 {
			ifo, err := os.Lstat(basePath)
			if err != nil {
				return err
			}

			if !ifo.IsDir() {
				basePath = filepath.Dir(basePath)
			}
		}
		basePath = filepath.Clean(basePath)
		if strings.HasSuffix(basePath, string(os.PathSeparator)) == false {
			basePath = basePath + string(os.PathSeparator)
		}
	}

	// 收集压缩的文件
	var fileList []string
	for _, file := range newFiles {
		err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			fileList = append(fileList, path)
			return nil
		})

		if err != nil {
			return err
		}
	}

	gw := gzip.NewWriter(d)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	var compress = func(fileFullPath string) error {
		f, err := os.Open(fileFullPath)
		if err != nil {
			return err
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Format = tar.FormatGNU
		header.Name = strings.TrimPrefix(filepath.Clean(fileFullPath), basePath)
		log.Println(">>>>> compress >>>>>", header.Name)

		if err = tw.WriteHeader(header); err != nil {
			return err
		}

		progressReader.SetReader(f)
		if _, err = io.Copy(tw, progressReader); err != nil {
			return err
		}

		return nil
	}

	// 压缩文件
	for _, path := range fileList {
		if err := compress(path); err != nil {
			return err
		}
	}

	return nil
}

// GzipDeCompress 解压 tar.gz
func GzipDeCompress(tarFile, dest string, progress Progress) error {
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

	objDecoder := mahonia.NewDecoder("GBK")

	var decompress = func(hdr *tar.Header) error {
		utf8String := hdr.Name
		if objDecoder != nil {
			utf8String = objDecoder.ConvertString(hdr.Name)
		}
		filename := fmt.Sprintf("%v/%v", dest, utf8String)
		isDir := hdr.Typeflag == tar.TypeDir

		file, err := createFile(filename, isDir)
		if err != nil {
			return err
		}
		if !isDir {
			defer file.Close()

			_, err = io.Copy(file, tr)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		if err := decompress(hdr); err != nil {
			return err
		}
	}
	return nil
}
