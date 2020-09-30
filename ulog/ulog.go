package ulog


import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func init() {
	logs.GetBeeLogger().EnableFuncCallDepth(true)
	logs.GetBeeLogger().SetLogFuncCallDepth(3)

	_, filename := filepath.Split(os.Args[0])
	appName := strings.TrimRight(filename, filepath.Ext(filename))

	_ = os.Mkdir("log", os.ModePerm)
	logFilename := fmt.Sprintf("./log/%v_%v.log", appName, time.Now().Format("0102_150405"))

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

func Info(format string, v ...interface{}) {
	logs.GetBeeLogger().Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	logs.GetBeeLogger().Debug(format, v...)
}

func Error(format string, v ...interface{}) {
	stack := getFileLine(5)
	var lines = "\n"
	for i := 0; i < len(stack); i++ {
		lines += "  " + stack[i] + "\n"
	}

	s := fmt.Sprintf(format, v...)
	s += lines
	logs.GetBeeLogger().Error(s)
}

func Warn(format string, v ...interface{}) {
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
