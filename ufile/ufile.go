package ufile

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Md5File 文件hash
func Md5File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	w := md5.New()
	_, err = io.Copy(w, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", w.Sum(nil)), nil
}

// SHA1File 文件hash
func SHA1File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// SHA256File 文件hash
func SHA256File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha256.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// IsExists 判断所给路径文件/文件夹是否存在
func IsExists(path string) bool {
	_, err := os.Lstat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}

	return true
}

func ListDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

// Split 分割路径
func Split(fullPath string) (dir, file, ext string) {
	var tmpFile string

	dir, tmpFile = filepath.Split(fullPath)
	ext = filepath.Ext(tmpFile)

	index := strings.LastIndex(tmpFile, ext)
	if index >= 0 {
		file = tmpFile[:index]
	} else {
		file = tmpFile
	}

	return
}

// https://github.com/mattetti/filebuffer

// FileInfo 文件信息
func FileInfo(path string) (os.FileInfo, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := fp.Stat()
	if err != nil {
		_ = fp.Close()
		return nil, err
	}
	
	_ = fp.Close()

	return info, nil
}
