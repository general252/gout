package ulog

import (
	"time"

	"github.com/astaxie/beego/logs"
)

// newLogWrapper return a LoggerInterface
func newLogWrapper() logs.Logger {
	return innerObjectLog
}

var (
	innerObjectLog = innerNewLogWrapper()
)

func innerNewLogWrapper() *logWrapper {
	return &logWrapper{}
}

const (
	AdapterLogWrapper = "uLogWrapper"
)

type LogHandler func(when time.Time, msg string, level LogLevel)

type logWrapper struct {
	handle LogHandler
}

func (l *logWrapper) Init(config string) error {
	return nil
}

func (l *logWrapper) WriteMsg(when time.Time, msg string, level int) error {
	if l.handle != nil {
		l.handle(when, msg, LogLevel(level))
	}
	return nil
}

func (l *logWrapper) Destroy() {
}

func (l *logWrapper) Flush() {
}

func (l *logWrapper) SetLogHandle(handle LogHandler) {
	l.handle = handle
}

// SetLogHandle 设置log handle
func SetLogHandle(handle LogHandler) {
	innerObjectLog.SetLogHandle(handle)
}

func init() {
	logs.Register(AdapterLogWrapper, newLogWrapper)
}
