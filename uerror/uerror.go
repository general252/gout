package uerror

import (
	"fmt"

	"github.com/general252/gout/ucode"
	"github.com/general252/gout/ustack"
)

type uError struct {
	code      ucode.Code
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

func (c *uError) Err() error {
	if c == nil {
		return nil
	}
	return c.err
}
func (c *uError) appendMessage(msg string) *uError {
	c.msg = append(c.msg, msg)
	return c
}
func (c *uError) Message() []string {
	if c == nil {
		return nil
	}
	return c.msg
}

func (c *uError) Msg() string {
	if c == nil || c.msg == nil {
		return ""
	}

	return mFormatString(c.msg)
}

func (c *uError) Stack() []string {
	if c == nil {
		return nil
	}
	return c.callstack
}

func (c *uError) WithError(err error) *uError {
	if c != nil {
		c.err = err
	}

	return c
}

func (c *uError) WithCode(code ucode.Code) *uError {
	if c != nil {
		c.code = code
	}

	return c
}

// AsUError convert to uError
func AsUError(err error) (*uError, bool) {
	uErr, ok := err.(*uError)
	return uErr, ok
}

// AsUErr convert to uError
func AsUErr(err error) *uError {
	uErr, _ := err.(*uError)
	return uErr
}

func GetOriError(err error) error {
	if ue := AsUErr(err); ue != nil {
		return ue.err
	}

	return err
}

func GetOriCode(err error) ucode.Code {
	if ue := AsUErr(err); ue != nil {
		return ue.code
	}

	return ucode.Unknown
}

func newUError(err error, msg string) *uError {
	if uErr, ok := AsUError(err); ok {
		return uErr.appendMessage(msg)
	}

	return &uError{
		callstack: ustack.CallStackList(3, 4),
		err:       err,
		msg:       []string{msg},
		code:      ucode.Unknown,
	}
}

// WithMessage fmt.Println(...)
func WithMessage(a ...interface{}) *uError {
	return newUError(nil, fmt.Sprint(a...))
}

// WithMessageF fmt.Printf(...)
func WithMessageF(format string, a ...interface{}) *uError {
	return newUError(nil, fmt.Sprintf(format, a...))
}

// WithNothing no msg, no error, only stack
func WithNothing() *uError {
	return newUError(nil, "")
}

// WithError no message
func WithError(err error) *uError {
	return newUError(err, "")
}

// WithErrorAndMessage fmt.Println(...) with error
func WithErrorAndMessage(err error, a ...interface{}) *uError {
	return newUError(err, fmt.Sprint(a...))
}

// WithErrorAndMessageF fmt.Printf(...) with error
func WithErrorAndMessageF(err error, format string, a ...interface{}) *uError {
	return newUError(err, fmt.Sprintf(format, a...))
}

func mFormatString(lines []string) string {
	var rs = ""
	for i := 0; i < len(lines); i++ {
		rs += fmt.Sprintf("  %02v. %v\n", i+1, lines[i])
	}

	return rs
}
