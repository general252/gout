package uerror

import (
	"fmt"
	"github.com/general252/gout/usafe"
)

type uError struct {
	err       error
	msg       []string
	callstack []string
}

// Error interface error
func (c *uError) Error() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf(
		"error: %v\nmessage: \n%vstack: \n%v",
		c.err, mFormatString(c.msg), mFormatString(c.callstack),
	)
}

func (c *uError) GetError() error {
	if c == nil {
		return nil
	}
	return c.err
}
func (c *uError) appendMessage(msg string) *uError {
	c.msg = append(c.msg, msg)
	return c
}
func (c *uError) GetMessage() []string {
	if c == nil {
		return nil
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
	if uErr, ok := ConvertToUError(err); ok {
		return uErr.appendMessage(msg)
	}

	return &uError{
		callstack: usafe.CallStackList(3, 4),
		err:       err,
		msg:       []string{msg},
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

func mFormatString(lines []string) string {
	var rs = ""
	for i := 0; i < len(lines); i++ {
		rs += fmt.Sprintf("  %02v. %v\n", i+1, lines[i])
	}

	return rs
}
