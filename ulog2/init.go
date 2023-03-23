package ulog2

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

var xLogSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	_, f, ext := split(file)
	xLogSourceDir = strings.ReplaceAll(file, fmt.Sprintf("%v%v", f, ext), "")
}

type StackFrame struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

// GetLastCallStack 获取当前代码位置
func GetLastCallStack() StackFrame {
	var shortFile = true
	shortFile = true

	for i := 0; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			if strings.HasSuffix(file, "/runtime/proc.go") {
				break
			}

			if !strings.HasPrefix(file, xLogSourceDir) || strings.HasSuffix(file, "_test.go") {

				if shortFile {
					short := file
					for i := len(file) - 1; i > 0; i-- {
						if file[i] == '/' {
							short = file[i+1:]
							break
						}
					}
					file = short
				}

				return StackFrame{
					File: file,
					Line: line,
				}
			}
		}
	}

	return StackFrame{
		File: "???",
		Line: 0,
	}
}

// split 分割路径
func split(fullPath string) (dir, file, ext string) {
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

// Format 格式化
func Format(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}

	return innerFormat(v[0], v[1:]...)
}

// innerFormat 引用 github.com\astaxie\beego@v1.12.3\logs\log.go
func innerFormat(f interface{}, v ...interface{}) string {
	var msg string
	switch f := f.(type) {
	case string:
		msg = f
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}
