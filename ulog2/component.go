package ulog2

import (
	"fmt"
	"os"
	"time"
)

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})

	// AddTag 添加永久tag
	AddTag(tags ...string) Logger
	// WithTag 添加临时tag
	WithTag(tags ...string) Logger
	GetTag() Tag
	HasTag(tag string) bool

	WithWriter(w LoggerWriter)
}

// SetDefaultWriter set default log io writer
func SetDefaultWriter(f LoggerWriter) {
	defaultLogWriter = f
}

// Component component
func Component(tags ...string) Logger {
	c := &componentData{
		tag:    tags,
		writer: defaultLogWriter,
	}

	return c
}

type componentData struct {
	tag    Tag
	tagTmp Tag
	writer LoggerWriter
}

func (tis *componentData) resetTempTag() Tag {
	tag := Tags(tis.tag...)
	tag.Append(tis.tagTmp...)

	tis.tagTmp = nil

	return tag
}

func (tis *componentData) Debug(v ...interface{}) {
	tis.writer(&JsonLogObject{
		Format:   logFormat,
		Time:     time.Now(),
		Location: getLastCallStack(),
		Level:    LevelDebug,
		Tag:      tis.resetTempTag(),
		Data:     Format(v...),
	})
}

func (tis *componentData) Info(v ...interface{}) {
	tis.writer(&JsonLogObject{
		Format:   logFormat,
		Time:     time.Now(),
		Location: getLastCallStack(),
		Level:    LevelInfo,
		Tag:      tis.resetTempTag(),
		Data:     Format(v...),
	})
}

func (tis *componentData) Warn(v ...interface{}) {
	tis.writer(&JsonLogObject{
		Format:   logFormat,
		Time:     time.Now(),
		Location: getLastCallStack(),
		Level:    LevelWarn,
		Tag:      tis.resetTempTag(),
		Data:     Format(v...),
	})
}

func (tis *componentData) Error(v ...interface{}) {
	tis.writer(&JsonLogObject{
		Format:   logFormat,
		Time:     time.Now(),
		Location: getLastCallStack(),
		Level:    LevelError,
		Tag:      tis.resetTempTag(),
		Data:     Format(v...),
	})
}

func (tis *componentData) AddTag(tags ...string) Logger {
	tis.tag.Append(tags...)
	return tis
}

func (tis *componentData) WithTag(tags ...string) Logger {
	tis.tagTmp.Append(tags...)
	return tis
}

func (tis *componentData) GetTag() Tag {
	return tis.tag
}

func (tis *componentData) HasTag(tag string) bool {
	for _, s := range tis.tag {
		if s == tag {
			return true
		}
	}
	return false
}

func (tis *componentData) WithWriter(writer LoggerWriter) {
	tis.writer = writer
}

const (
	logFormat = "%s %s:%d [%s] %s%s\n"
	dateTime  = "2006-01-02 15:04:05"
)

// JsonLogObject 日志信息
type JsonLogObject struct {
	Format   string     `json:"format"`   // 默认格式化 "%v %v:%v [%c] %v %v\n"
	Time     time.Time  `json:"time"`     // 时间
	Location StackFrame `json:"location"` // 日志调用位置
	Level    LogLevel   `json:"level"`    // 日志级别
	Tag      Tag        `json:"tag"`      // tag
	Data     string     `json:"data"`     // 日志数据
}

func (tis *JsonLogObject) String() string {
	return _getLogString(tis)
}

type LoggerWriter func(o *JsonLogObject)

var (
	_getLogString = func(o *JsonLogObject) string {
		return fmt.Sprintf(o.Format, o.Time.Format(dateTime), o.Location.File, o.Location.Line, o.Level, o.Tag.String(), o.Data)
	}

	defaultLogWriter = func(o *JsonLogObject) {
		_, _ = fmt.Fprintf(os.Stderr, _getLogString(o))
	}
)
