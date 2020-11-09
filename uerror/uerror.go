package uerror

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type uError struct {
	err       error
	msg       string
	callstack []string
}

// Error interface error
func (c *uError) Error() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("error: %v\nmessage: %v\nstack: \n%v",
		c.err, c.msg, mCallStackFormatString(c.callstack))
}

func (c *uError) GetError() error {
	if c == nil {
		return nil
	}
	return c.err
}

func (c *uError) GetMessage() string {
	if c == nil {
		return ""
	}
	return c.msg
}

func (c *uError) GetStack() []string {
	if c == nil {
		return nil
	}
	return c.callstack
}

// ConvertToUError convert to uError
func ConvertToUError(err error) (*uError, bool) {
	uErr, ok := err.(*uError)
	return uErr, ok
}

func newUError(err error, msg string) *uError {
	return &uError{
		callstack: CallStackList(3, 4),
		err:       err,
		msg:       msg,
	}
}

// Println fmt.Println(...)
func Println(a ...interface{}) *uError {
	return newUError(nil, fmt.Sprint(a...))
}

// Printf fmt.Printf(...)
func Printf(format string, a ...interface{}) *uError {
	return newUError(nil, fmt.Sprintf(format, a...))
}

// PrintlnWithError fmt.Println(...) with error
func PrintlnWithError(err error, a ...interface{}) *uError {
	return newUError(err, fmt.Sprint(a...))
}

// PrintfWithError fmt.Printf(...) with error
func PrintfWithError(err error, format string, a ...interface{}) *uError {
	return newUError(err, fmt.Sprintf(format, a...))
}

// CallStackList 获取调用堆栈, 0: 本函数, 1: 上级函数, 2: 上上级函数
// startDepth: 从n级统计,
// count: 统计count个
func CallStackList(startDepth, count int) []string {
	var lines []string
	for i := startDepth; i <= startDepth+count; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			file = "???"
			line = 0
			break
		} else {
			_, file = filepath.Split(file)
		}

		lines = append(lines, fmt.Sprintf("%v:%v", file, line))
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

func mCallStackFormatString(stack []string) string {
	var lines = ""
	for i := 0; i < len(stack); i++ {
		lines += fmt.Sprintf("  %02v. %v\n", i+1, stack[i])
	}

	return lines
}
