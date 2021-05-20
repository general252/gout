package ushell

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

// ShellCommand shell
// ShellCommand("ls -ltr")
// ShellCommand("SystemInfo")
func ShellCommand(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	var strOut = stdout.String()
	var strErr = stderr.String()

	var dec = mahonia.NewDecoder("gbk")

	strOut = dec.ConvertString(strOut)
	strErr = dec.ConvertString(strErr)

	return strOut, strErr, err
}

// ShellCommandStream("ls -ltr", func([]byte, bool){ })  isError: true是stderr输出, false: 是stdout输出
// ShellCommandStream("SystemInfo", func([]byte, bool){ })
func ShellCommandStream(command string, cb func(msg []byte, isError bool)) error {
	if cb == nil {
		return fmt.Errorf("cb is nil")
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		scanner := bufio.NewReader(stdout)
		for true {
			msg, _, err := scanner.ReadLine()
			if err != nil {
				return
			}
			cb(msg, false)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		scanner := bufio.NewReader(stderr)
		for true {
			msg, _, err := scanner.ReadLine()
			if err != nil {
				return
			}
			cb(msg, true)
		}
	}()

	if err := cmd.Run(); err != nil {
		return err
	}

	wg.Wait()

	return err
}

// OpenUri open uri on browser
func OpenUri(uri string) error {
	var commands = map[string]string{
		"windows": "cmd /c start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}

	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	cmdParams := strings.Split(run, " ")
	cmdParams = append(cmdParams, uri)

	cmd := exec.Command(cmdParams[0], cmdParams[1:]...)
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return cmd.Run()
}
