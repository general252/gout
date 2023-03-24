package ulog2

import (
	"log"
	"testing"
)

func a() {
	stacks := GetLastCallStackDepth(10)
	for _, stack := range stacks {
		log.Printf("%v:%v %v", stack.File, stack.Line, stack.Func)
	}
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
