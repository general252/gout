package uapp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

func GetExePath() (string, error) {
	appFilePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}

	return appFilePath, nil
}

func GetExeName() string {
	_, filename := filepath.Split(os.Args[0])
	ext := filepath.Ext(os.Args[0])
	index := strings.LastIndex(filename, ext)
	if index == -1 {
		return filename
	}

	return filename[:index]
}

func GetExeDir() (string, error) {
	dir := filepath.Dir(os.Args[0])
	absDir, err := filepath.Abs(dir)

	return absDir, err
}

func UnHandlerException() {
	errs := recover()
	if errs == nil {
		return
	}

	fileName := fmt.Sprintf("%s_%s.dmp",
		GetExeName(),
		time.Now().Format("20060102_150405")) //保存错误信息文件名:程序名-当前时间（年月日时分秒）

	f, err := os.Create(fileName)
	if err == nil {
		_, _ = f.WriteString(fmt.Sprintf("%v\r\n", errs)) //输出panic信息
		_, _ = f.WriteString("========\r\n")

		_, _ = f.WriteString(string(debug.Stack())) //输出堆栈信息

		_ = f.Close()
	}
}

func ReStartApp() {
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("restart error %v\n", err)
	} else {
		fmt.Println("success")
	}

	os.Exit(0)
}
