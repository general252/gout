package ulog2

import (
	"fmt"
	"log"
	"testing"
)

func a() {
	stacks := getLastCallStackDepth(6)
	for _, stack := range stacks {
		log.Printf("%v:%v %v", stack.File, stack.Line, stack.Func)
	}
	fmt.Println(stacks.String())
}

func b() {
	a()
}

func c() {
	b()
}

func TestGetLastCallStackDepth(t *testing.T) {
	c()
}
