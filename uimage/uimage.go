package uimage

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"

	_ "golang.org/x/image/bmp"
)

// Decode 解析图像(支持bmp/gig/jpeg/png)
func Decode(imgFilePath string) (image.Image, string, error) {
	fp, err := os.Open(imgFilePath)
	if err != nil {
		return nil, "", err
	}
	defer fp.Close()

	return image.Decode(fp)
}

// DecodeBytes 解析图像(支持bmp/gig/jpeg/png)
func DecodeBytes(imageData []byte) (image.Image, string, error) {
	buffer := bytes.NewBuffer(imageData)

	return image.Decode(buffer)
}

// DecodeConfig x
func DecodeConfig(imgFilePath string) (image.Config, string, error) {
	fp, err := os.Open(imgFilePath)
	if err != nil {
		return image.Config{}, "", err
	}
	defer fp.Close()

	return image.DecodeConfig(fp)
}

// DecodeConfigBytes x
func DecodeConfigBytes(imageData []byte) (image.Config, string, error) {
	buffer := bytes.NewBuffer(imageData)

	return image.DecodeConfig(buffer)
}

// DetectImageType http content type.  default "application/octet-stream", ["image/jpeg", "image/gif", "image/png", "image/webp"]
func DetectImageType(imageData []byte) {
	http.DetectContentType(imageData)
}

/*
var GIF = []byte("GIF")
var BMP = []byte("BM")
var JPG = []byte{0xff, 0xd8, 0xff}
var PNG = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}

type ImageType int

const (
	UnKnowType ImageType = iota
	GifType
	BmpType
	JpgType
	PngType
)

func GetImageType(imgFilePath string) (ImageType, error) {
	f, err := os.Open(imgFilePath)
	if err != nil {
		return UnKnowType, err
	}
	defer f.Close()

	buffer := make([]byte, 16)
	if _, err := f.Read(buffer); err != nil {
		return UnKnowType, err
	}

	if bytes.Equal(PNG, buffer[0:8]) {
		return PngType, nil
	}
	if bytes.Equal(GIF, buffer[0:3]) {
		return GifType, nil
	}
	if bytes.Equal(BMP, buffer[0:2]) {
		return BmpType, nil
	}
	if bytes.Equal(JPG, buffer[0:3]) {
		return JpgType, nil
	}

	return UnKnowType, fmt.Errorf("undefined type")
}
*/
