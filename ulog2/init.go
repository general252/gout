package ulog2

import (
	"fmt"
	"runtime"
	"strings"
)

var xLogSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	_, f, ext := split(file)
	xLogSourceDir = strings.ReplaceAll(file, fmt.Sprintf("%v%v", f, ext), "")
}
