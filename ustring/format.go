package ustring

import (
	"fmt"
	"strings"
)

// Format 格式化
func Format(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}

	return innerFormat(v[0], v[1:]...)
}

// innerFormat 引用 github.com\astaxie\beego@v1.12.3\logs\log.go
func innerFormat(f interface{}, v ...interface{}) string {
	var msg string
	switch f := f.(type) {
	case string:
		msg = f
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}
