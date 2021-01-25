package ulog

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/general252/gout/uapp"
	"github.com/general252/gout/ufile"
	"github.com/general252/gout/ustack"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"
)

type LogLevel int

// RFC5424 log message levels.
const (
	LevelEmergency LogLevel = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)
const (
	LevelInfo  = LevelInformational
	LevelTrace = LevelDebug
	LevelWarn  = LevelWarning
)

// FileConfig 文件配置
type FileConfig struct {
	Filename string `json:"filename,omitempty"` // filename 保存的文件名

	// Rotate at line
	MaxLines int `json:"maxlines,omitempty"` // maxlines 每个文件保存的最大行数，默认值 1000000

	MaxFiles int `json:"maxfiles,omitempty"`

	// Rotate at size
	MaxSize int `json:"maxsize,omitempty"` // maxsize 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB

	// Rotate daily
	Daily   bool  `json:"daily,omitempty"`   // 是否按照每天 logrotate，默认是 true
	MaxDays int64 `json:"maxdays,omitempty"` // 文件最多保存多少天，默认保存 7 天

	// Rotate hourly
	Hourly   bool  `json:"hourly,omitempty"`
	MaxHours int64 `json:"maxhours,omitempty"`

	Rotate     bool     `json:"rotate,omitempty"` // 是否开启 logrotate，默认是 true
	Level      LogLevel `json:"level,omitempty"`  // 日志保存的时候的级别，默认是 Trace 级别
	Perm       string   `json:"perm,omitempty"`   // 日志文件权限
	RotatePerm string   `json:"rotateperm,omitempty"`
}

// ConsoleConfig 控制台配置
type ConsoleConfig struct {
	Level    LogLevel `json:"level,omitempty"`
	Colorful bool     `json:"color,omitempty"` //this filed is useful only when system's terminal supports color
}

type HandleError func(err string)

var defaultHandleError HandleError = func(err string) {
}

// SetDefaultHandleError 设置默认错误处理
func SetDefaultHandleError(handleRecover HandleError) {
	defaultHandleError = handleRecover
}

func initLogPath() string {
	// 创建日志
	appDir, _ := uapp.GetExeDir()
	logDir := fmt.Sprintf("%v/log", appDir)

	if !ufile.IsExists(logDir) {
		_ = os.Mkdir("log", os.ModePerm)
	}

	return logDir
}

func init() {
	logs.GetBeeLogger().EnableFuncCallDepth(true)
	logs.GetBeeLogger().SetLogFuncCallDepth(3)

	logFilename := fmt.Sprintf("%v/%v_%v_pid%v.log",
		initLogPath(),
		uapp.GetExeName(),
		time.Now().Format("20060102_150405"),
		os.Getpid())

	// 文件
	if true {
		fileConfigStr, _ := json.Marshal(&FileConfig{
			Filename: logFilename,
			MaxFiles: 100,
			MaxSize:  1 << 24,
			MaxDays:  15,
			Level:    LevelTrace,
		})
		if err := logs.SetLogger(logs.AdapterFile, string(fileConfigStr)); err != nil {
			log.Printf("set logger fail %v %v\n", logs.AdapterFile, err)
		}
	}

	// 控制台
	if true {
		consoleConfigStr, _ := json.Marshal(&ConsoleConfig{
			Level:    LevelTrace,
			Colorful: true,
		})
		if err := logs.SetLogger(logs.AdapterConsole, string(consoleConfigStr)); err != nil {
			log.Printf("set logger fail %v %v\n", logs.AdapterConsole, err)
		}
	}

	// 自定义
	if true {
		if err := logs.SetLogger(AdapterLogWrapper, "{}"); err != nil {
			log.Printf("set logger fail %v %v\n", AdapterLogWrapper, err)
		}
	}

	logs.GetBeeLogger().Async(1000)
}

// SetFileLogConfig 设置日志文件配置
func SetFileLogConfig(enable bool, config FileConfig) error {
	if enable == false {
		err := logs.GetBeeLogger().DelLogger(logs.AdapterFile)
		return err
	}

	bytesConfig, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// 删除旧配置(如果有)
	_ = logs.GetBeeLogger().DelLogger(logs.AdapterFile)

	// 创建新配置
	if err = logs.SetLogger(logs.AdapterFile, string(bytesConfig)); err != nil {
		return err
	}

	return nil
}

// SetConsoleConfig 设置控制台配置
func SetConsoleConfig(enable bool, config ConsoleConfig) error {
	if enable == false {
		err := logs.GetBeeLogger().DelLogger(logs.AdapterConsole)
		return err
	}

	bytesConfig, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// 删除旧配置(如果有)
	_ = logs.GetBeeLogger().DelLogger(logs.AdapterConsole)

	// 创建新配置
	if err = logs.SetLogger(logs.AdapterConsole, string(bytesConfig)); err != nil {
		return err
	}

	return nil
}

// TaskCheckLog 限制日志总大小
func TaskCheckLog(limitMB int64) {
	fl := GetLogFiles()

	limitBytes := limitMB * 1024 * 1024

	var totalSize int64 = 0
	for i := 0; i < len(fl); i++ {
		var info = fl[i]
		totalSize += info.Size
		if totalSize > limitBytes {
			delFile := fl[i:]
			for _, fileInfo := range delFile {
				log.Printf("remove log file %v", fileInfo.Path)
				_ = os.Remove(fileInfo.Path)
			}

			break
		}
	}
}

type LogFileInfo struct {
	Path    string
	ModTime time.Time
	Size    int64
}

// GetLogFiles 获取日志文件
func GetLogFiles() []LogFileInfo {
	logRootPath := initLogPath()
	fl, err := ufile.ListDir(logRootPath)
	if err != nil {
		return nil
	}

	var result []LogFileInfo
	for _, info := range fl {
		if info.IsDir() {
			continue
		}
		if filepath.Ext(info.Name()) != ".log" {
			continue
		}

		path := filepath.Join(logRootPath, info.Name())
		result = append(result, LogFileInfo{
			Path:    path,
			ModTime: info.ModTime(),
			Size:    info.Size(),
		})
	}
	if result != nil {
		sort.Slice(result, func(i, j int) bool {
			return result[i].ModTime.Unix() > result[j].ModTime.Unix()
		})
	}

	return result
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
