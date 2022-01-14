package ucolor

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestYellow(t *testing.T) {
	fmt.Printf("hello %v\n", Black("Black"))
	fmt.Printf("hello %v\n", Red("Red"))
	fmt.Printf("hello %v\n", Green("Green"))
	fmt.Printf("hello %v\n", Yellow("Yellow"))
	fmt.Printf("hello %v\n", Blue("Blue"))
	fmt.Printf("hello %v\n", Magenta("Magenta"))
	fmt.Printf("hello %v\n", Cyan("Cyan"))
	fmt.Printf("hello %v\n", White("White"))

	fmt.Printf("hello %v\n", HiBlack("HiBlack"))
	fmt.Printf("hello %v\n", HiRed("HiRed"))
	fmt.Printf("hello %v\n", HiGreen("HiGreen"))
	fmt.Printf("hello %v\n", HiYellow("HiYellow"))
	fmt.Printf("hello %v\n", HiBlue("HiBlue"))
	fmt.Printf("hello %v\n", HiMagenta("HiMagenta"))
	fmt.Printf("hello %v\n", HiCyan("HiCyan"))
	fmt.Printf("hello %v\n", HiWhite("HiWhite"))
}

func TestRed(t *testing.T) {
	var output io.Writer = os.Stderr
	// output = colorable.NewColorableStderr() // "github.com/mattn/go-colorable"

	_, _ = fmt.Fprintf(output, "hello %v", Red("Red"))
}

func TestColorStringX(t *testing.T) {
	s := ColorStringX([]FontColor{Bold, FgHiCyan, BgWhite, Underline}, "abc")

	fmt.Printf("%s", s)
}
