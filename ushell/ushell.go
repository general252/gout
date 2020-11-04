package ushell

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/general252/gout/uhost"
	"os/exec"
	"runtime"
	"strings"
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
