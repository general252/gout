package ustack

import (
	"fmt"
	"path/filepath"
	"runtime"
)

var (
	addFunction = false
)

func CallStackAddFunction(haveFunction bool) {
	addFunction = haveFunction
}

// CallStackList 获取调用堆栈, 0: 本函数, 1: 上级函数, 2: 上上级函数
// startDepth: 从n级统计,
// count: 统计count个
func CallStackList(startDepth, count int) []string {
	var lines []string
	for i := startDepth; i <= startDepth+count; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			file = "???"
			line = 0
			break
		} else {
			_, file = filepath.Split(file)
		}

		var newLine = fmt.Sprintf("%v:%v", file, line)
		if addFunction {
			if fn := runtime.FuncForPC(pc); fn != nil {
				newLine = fmt.Sprintf("%-20s func: %v(..)", newLine, fn.Name())
			}
		}

		lines = append(lines, newLine)
	}

	return lines
}

// CallStackFormatString 获取调用堆栈
func CallStackFormatString(startDepth, count int) string {
	stack := CallStackList(startDepth, count)
	var lines = "\n"
	for i := 0; i < len(stack); i++ {
		lines += "  " + stack[i] + "\n"
	}

	return lines
}

// RoutineCount routine个数
func RoutineCount() int {
	return runtime.NumGoroutine()
}

// RoutineStacks 所有routine的调用堆栈
func RoutineStacks() []byte {
	// We don't know how big the buffer needs to be to collect
	// all the goroutines. Start with 1 MB and try a few times, doubling each time.
	// Give up and use a truncated trace if 64 MB is not enough.
	buf := make([]byte, 1<<20)
	for i := 0; ; i++ {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		if len(buf) >= 64<<20 {
			// Filled 64 MB - stop there.
			break
		}
		buf = make([]byte, 2*len(buf))
	}

	return buf
}
