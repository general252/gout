package ucolor

import (
	"fmt"
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

func ColorString(s string, c FontColor) string {
	return fmt.Sprintf("%s[%dm%s%s[%dm", Escape, c, s, Escape, Reset)
}

func Black(s string) string {
	return ColorString(s, FgBlack)
}
func Red(s string) string {
	return ColorString(s, FgRed)
}
func Green(s string) string {
	return ColorString(s, FgGreen)
}
func Yellow(s string) string {
	return ColorString(s, FgYellow)
}
func Blue(s string) string {
	return ColorString(s, FgBlue)
}
func Magenta(s string) string {
	return ColorString(s, FgMagenta)
}
func Cyan(s string) string {
	return ColorString(s, FgCyan)
}
func White(s string) string {
	return ColorString(s, FgWhite)
}

func HiBlack(s string) string {
	return ColorString(s, FgHiBlack)
}
func HiRed(s string) string {
	return ColorString(s, FgHiRed)
}
func HiGreen(s string) string {
	return ColorString(s, FgHiGreen)
}
func HiYellow(s string) string {
	return ColorString(s, FgHiYellow)
}
func HiBlue(s string) string {
	return ColorString(s, FgHiBlue)
}
func HiMagenta(s string) string {
	return ColorString(s, FgHiMagenta)
}
func HiCyan(s string) string {
	return ColorString(s, FgHiCyan)
}
func HiWhite(s string) string {
	return ColorString(s, FgHiWhite)
}
