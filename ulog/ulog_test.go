package ulog

import (
	"fmt"
	"github.com/general252/gout/utask"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ExampleSetFileLogConfig() {
	_, filename := filepath.Split(os.Args[0])
	appName := strings.TrimRight(filename, filepath.Ext(filename))

	logFilename := fmt.Sprintf("./log/%v_%v_%v.log", appName, time.Now().Format("20060102_150405"), os.Getpid())
	var config = FileConfig{
		Filename:   logFilename,
		Daily:      true,
		MaxDays:    7,
		Hourly:     false,
		MaxHours:   168,
		Rotate:     true,
		RotatePerm: "0440",
		Level:      LevelTrace,
		Perm:       "0660",
		MaxLines:   10000000,
		MaxFiles:   999,
		MaxSize:    1 << 28,
	}

	if err := SetFileLogConfig(config); err != nil {
		fmt.Printf("set file log config fail %v", err)
	}

	// output:
	//
}

func ExampleTaskCheckLog() {
	var check = func() {
		TaskCheckLog(1024)
	}
	_, _ = utask.AddTaskEvery(time.Minute, check)
}
