package uzip

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ZipCompress(files []string, dest string, progress Progress) error {
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

			fileList = append(fileList, path)
			return nil
		})

		if err != nil {
			return err
		}
	}

	archive := zip.NewWriter(d)
	defer archive.Close()

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

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(fileFullPath, basePath)
		log.Println(">>>>> compress >>>>>", header.Name)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			progressReader.SetReader(f)
			if _, err = io.Copy(writer, progressReader); err != nil {
				return err
			}
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

func ZipDeCompress(tarFile, dest string, progress Progress) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	progressReader := NewProcessReader(progress)
	info, err := srcFile.Stat()
	if err != nil {
		return err
	}
	progressReader.totalSize = info.Size()
	progressReader.SetReaderAt(srcFile)

	zipReader, err := zip.NewReader(progressReader, info.Size())
	if err != nil {
		return err
	}

	var decompress = func(f *zip.File) error {
		filename := filepath.Join(dest, f.Name)
		isDir := f.FileInfo().IsDir()

		file, err := createFile(filename, isDir)
		if err != nil {
			return err
		}

		if !isDir {
			defer file.Close()

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			_, err = io.Copy(file, inFile)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, f := range zipReader.File {
		if err := decompress(f); err != nil {
			return err
		}
	}

	return nil
}
