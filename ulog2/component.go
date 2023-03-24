package ulog2

import (
	"fmt"
	"os"
	"sync"
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
	WithStack(depth int) Logger
}

// SetDefaultWriter set default log io writer
func SetDefaultWriter(f LoggerWriter) {
	defaultLogWriter = f
}

// Component component
func Component(tags ...string) Logger {
	c := &componentData{
		tag:        tags,
		writer:     defaultLogWriter,
		withStack:  false,
		stackDepth: 4,
	}

	return c
}

type componentData struct {
	tag    Tag
	tagTmp Tag
	writer LoggerWriter

	withStack  bool
	stackDepth int

	mux sync.RWMutex
}

func (tis *componentData) WithStack(depth int) Logger {
	return &componentData{
		tag:        tis.tag,
		tagTmp:     tis.tagTmp,
		writer:     tis.writer,
		withStack:  true,
		stackDepth: depth,
	}
}

func (tis *componentData) WithWriter(writer LoggerWriter) {
	tis.writer = writer
}

func (tis *componentData) Debug(v ...interface{}) {
	tis.print(LevelDebug, v...)
}

func (tis *componentData) Info(v ...interface{}) {
	tis.print(LevelInfo, v...)
}

func (tis *componentData) Warn(v ...interface{}) {
	tis.print(LevelWarn, v...)
}

func (tis *componentData) Error(v ...interface{}) {
	tis.print(LevelError, v...)
}

func (tis *componentData) AddTag(tags ...string) Logger {
	tis.mux.Lock()
	defer tis.mux.Unlock()

	for _, tag := range tags {
		find := false
		for _, s := range tis.tag {
			if s == tag {
				find = true
				break
			}
		}

		if !find {
			tis.tag.Append(tag)
		}
	}
	return tis
}

func (tis *componentData) WithTag(tags ...string) Logger {
	object := &componentData{
		tag:        tis.tag,
		tagTmp:     tis.tagTmp,
		writer:     tis.writer,
		withStack:  tis.withStack,
		stackDepth: tis.stackDepth,
	}

	object.tagTmp.Append(tags...)

	return object
}

func (tis *componentData) GetTag() Tag {
	tis.mux.RLock()
	defer tis.mux.RUnlock()

	return tis.tag
}

func (tis *componentData) HasTag(tag string) bool {
	tis.mux.RLock()
	defer tis.mux.RUnlock()

	for _, s := range tis.tag {
		if s == tag {
			return true
		}
	}
	return false
}

func (tis *componentData) getTag() Tag {
	tis.mux.RLock()
	defer tis.mux.RUnlock()

	if tis.tagTmp == nil {
		return tis.tag
	}

	tag := Tags(tis.tag...)
	tag.Append(tis.tagTmp...)

	return tag
}

func (tis *componentData) getStackDepth() int {
	if tis.stackDepth < 1 {
		return 1
	}
	return tis.stackDepth
}

func (tis *componentData) reset() {
	tis.mux.Lock()
	defer tis.mux.Unlock()

	tis.tagTmp = nil
	tis.withStack = false
	tis.stackDepth = 4
}

var eventPool = &sync.Pool{
	New: func() interface{} {
		return &JsonLogObject{}
	},
}

func (tis *componentData) print(level LogLevel, v ...interface{}) {
	var logData *JsonLogObject
	if true {
		logData = eventPool.Get().(*JsonLogObject)
		defer eventPool.Put(logData)

		logData.Time = time.Now()
		logData.Level = level
		logData.Tag = tis.getTag()
		logData.Data = Format(v...)
	} else {
		logData = &JsonLogObject{
			Time:  time.Now(),
			Level: level,
			Tag:   tis.getTag(),
			Data:  Format(v...),
		}
	}

	if tis.withStack {
		logData.Stacks = getLastCallStackDepth(tis.getStackDepth())
		logData.Location = logData.Stacks[0]
	} else {
		logData.Location = getLastCallStackDepth(1)[0]
	}

	tis.reset()

	tis.writer(logData)
}

const (
	dateTime = "2006-01-02 15:04:05"
)

// JsonLogObject 日志信息
type JsonLogObject struct {
	Time     time.Time  `json:"time"`     // 时间
	Location StackFrame `json:"location"` // 日志调用位置
	Level    LogLevel   `json:"level"`    // 日志级别
	Tag      Tag        `json:"tag"`      // tag
	Data     string     `json:"data"`     // 日志数据
	Stacks   Stacks     `json:"stacks"`   // 调用堆栈
}

func (tis *JsonLogObject) String() string {
	return _getLogString(tis)
}

type LoggerWriter func(o *JsonLogObject)

const (
	logFormat      = "%s %s:%d [%s] %s%s\n"
	logFormatStack = "%s %s:%d [%s] %s%s ↘↘↘\n%s\n"
)

var (
	_getLogString = func(o *JsonLogObject) string {
		if o.Stacks != nil {
			return fmt.Sprintf(logFormatStack, o.Time.Format(dateTime), o.Location.File, o.Location.Line, o.Level, o.Tag.String(), o.Data, o.Stacks.String())
		} else {
			return fmt.Sprintf(logFormat, o.Time.Format(dateTime), o.Location.File, o.Location.Line, o.Level, o.Tag.String(), o.Data)
		}
	}

	defaultLogWriter = func(o *JsonLogObject) {
		_, _ = fmt.Fprintf(os.Stderr, _getLogString(o))
	}
)
