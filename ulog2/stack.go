package ulog2

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Stacks []StackFrame

func (c *Stacks) String() string {
	if c == nil {
		return ""
	}
	var stacks = []StackFrame(*c)
	var count = len(stacks)
	if count > 0 {
		var r = "stack: [\n"
		for i := 0; i < count-1; i++ {
			r += fmt.Sprintf("  %v:%v %v\n", stacks[i].File, stacks[i].Line, stacks[i].Func)
		}
		r += fmt.Sprintf("  %v:%v %v\n]", stacks[count-1].File, stacks[count-1].Line, stacks[count-1].Func)

		return r
	}

	return ""
}

type StackFrame struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Func string `json:"func"`
}

// GetLastCallStackDepth 获取当前代码调用堆栈
func GetLastCallStackDepth(depth int) Stacks {
	var (
		shortFile = true
		result    Stacks
	)

	if depth <= 0 {
		return result
	}

	shortFile = true

	for i := 1; i < 15; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 判断调用位置
		if !strings.HasPrefix(file, xLogSourceDir) || strings.HasSuffix(file, "_test.go") {
			stack := StackFrame{
				File: file,
				Line: line,
			}

			// 简短文件名
			if shortFile {
				for i := len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						stack.File = file[i+1:]
						break
					}
				}
			}

			// 简短函数名
			if fn := runtime.FuncForPC(pc); fn != nil {
				name := fn.Name()
				stack.Func = name
				for i := len(name) - 1; i > 0; i-- {
					if name[i] == '.' {
						stack.Func = name[i+1:]
						break
					}
				}
			}

			result = append(result, stack)
			if len(result) == depth {
				break
			}
		}

		if strings.HasSuffix(file, "/runtime/proc.go") {
			break
		}
	}

	if len(result) == 0 {
		result = append(result, StackFrame{File: "???", Line: 0})
	}

	return result
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
