package ulog

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/general252/gout/ustack"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type HandleError func(err string)

var defaultHandleError HandleError = func(err string) {
}

// SetDefaultHandleError 设置默认错误处理
func SetDefaultHandleError(handleRecover HandleError) {
	defaultHandleError = handleRecover
}

func init() {
	logs.GetBeeLogger().EnableFuncCallDepth(true)
	logs.GetBeeLogger().SetLogFuncCallDepth(3)

	_, filename := filepath.Split(os.Args[0])
	appName := strings.TrimRight(filename, filepath.Ext(filename))

	_ = os.Mkdir("log", os.ModePerm)
	logFilename := fmt.Sprintf("./log/%v_%v_%v.log", appName, os.Getpid(), time.Now().Format("20060102_150405"))

	configStr := fmt.Sprintf(`{"level":%v, "filename":"%v", "maxdays":%v}`, logs.LevelDebug, logFilename, 15)
	if err := logs.SetLogger(logs.AdapterFile, configStr); err != nil {
	}

	configStr = fmt.Sprintf(`{"level":%v,"color":true}`, logs.LevelDebug)
	if err := logs.SetLogger(logs.AdapterConsole, configStr); err != nil {
	}

	logs.GetBeeLogger().Async(1000)
}

func Flush() {
	logs.GetBeeLogger().Flush()
}

func Info(v ...interface{}) {
	logs.GetBeeLogger().Info(fmt.Sprint(v...))
}
func Debug(v ...interface{}) {
	logs.GetBeeLogger().Debug(fmt.Sprint(v...))
}
func Error(v ...interface{}) {
	logs.GetBeeLogger().Error(fmt.Sprint(v...))
}
func Warn(v ...interface{}) {
	logs.GetBeeLogger().Warn(fmt.Sprint(v...))
}

func InfoF(format string, v ...interface{}) {
	logs.GetBeeLogger().Info(format, v...)
}

var (
	infoTimesMap = sync.Map{} // tag:times
	errorTimeMap = sync.Map{} // tag:times
)

// InfoFWithTimes 每times次输出一次日志
func InfoFWithTimes(tag string, times int, format string, v ...interface{}) {
	var valueTimes = 0
	objTimes, ok := infoTimesMap.Load(tag)
	if ok {
		valueTimes = objTimes.(int)
		infoTimesMap.Store(tag, valueTimes+1)
	} else {
		infoTimesMap.Store(tag, 1)
	}

	if valueTimes%times == 0 {
		InfoF(format, v...)
	}
}

// ErrorFWithTimes 每times次输出一次日志
func ErrorFWithTimes(tag string, times int, format string, v ...interface{}) {
	var valueTimes = 0
	objTimes, ok := errorTimeMap.Load(tag)
	if ok {
		valueTimes = objTimes.(int)
		errorTimeMap.Store(tag, valueTimes+1)
	} else {
		errorTimeMap.Store(tag, 1)
	}

	if valueTimes%times == 0 {
		ErrorF(format, v...)
	}
}

func DebugF(format string, v ...interface{}) {
	logs.GetBeeLogger().Debug(format, v...)
}

func ErrorF(format string, v ...interface{}) {
	//stack := getFileLine(5)
	stack := ustack.CallStackList(2, 1)
	var lines = "\nlog error stack:\n"
	for i := 0; i < len(stack); i++ {
		lines += fmt.Sprintf(" %2d. %v\n", i+1, stack[i])
	}

	s := fmt.Sprintf(format, v...)
	s += lines
	logs.GetBeeLogger().Error(s)
	defaultHandleError(s)
}

func WarnF(format string, v ...interface{}) {
	logs.GetBeeLogger().Warn(format, v...)
}

// 获取调用堆栈, 0: 本函数, 1: 上级函数, 2: 上上级函数
// callDepth: 多少个上级函数
func getFileLine(callDepth int) []string {
	var lines []string
	for i := 2; i <= callDepth; i++ {
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
