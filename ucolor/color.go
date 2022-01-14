package ucolor

import (
	"fmt"
	"strings"

	"github.com/general252/gout/ustring"
)

// FontColor defines a single SGR Code
type FontColor int

// Escape [30m  [0m
const Escape = "\x1b"

// Base attributes
const (
	Reset FontColor = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack FontColor = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack FontColor = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack FontColor = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack FontColor = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

func ColorString(c FontColor, v ...interface{}) string {
	return fmt.Sprintf("%s[%dm%s%s[%dm", Escape, c, ustring.Format(v...), Escape, Reset)
}

func ColorStringX(colors []FontColor, v ...interface{}) string {
	var colorArray []string
	for _, c := range colors {
		colorArray = append(colorArray, fmt.Sprintf("%d", c))
	}
	colorStr := strings.Join(colorArray, ";")

	return fmt.Sprintf("%s[%sm%s%s[%dm", Escape, colorStr, ustring.Format(v...), Escape, Reset)
}

func Black(v ...interface{}) string {
	return ColorString(FgBlack, v...)
}
func Red(v ...interface{}) string {
	return ColorString(FgRed, v...)
}
func Green(v ...interface{}) string {
	return ColorString(FgGreen, v...)
}
func Yellow(v ...interface{}) string {
	return ColorString(FgYellow, v...)
}
func Blue(v ...interface{}) string {
	return ColorString(FgBlue, v...)
}
func Magenta(v ...interface{}) string {
	return ColorString(FgMagenta, v...)
}
func Cyan(v ...interface{}) string {
	return ColorString(FgCyan, v...)
}
func White(v ...interface{}) string {
	return ColorString(FgWhite, v...)
}

func HiBlack(v ...interface{}) string {
	return ColorString(FgHiBlack, v...)
}
func HiRed(v ...interface{}) string {
	return ColorString(FgHiRed, v...)
}
func HiGreen(v ...interface{}) string {
	return ColorString(FgHiGreen, v...)
}
func HiYellow(v ...interface{}) string {
	return ColorString(FgHiYellow, v...)
}
func HiBlue(v ...interface{}) string {
	return ColorString(FgHiBlue, v...)
}
func HiMagenta(v ...interface{}) string {
	return ColorString(FgHiMagenta, v...)
}
func HiCyan(v ...interface{}) string {
	return ColorString(FgHiCyan, v...)
}
func HiWhite(v ...interface{}) string {
	return ColorString(FgHiWhite, v...)
}
