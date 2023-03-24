package ulog2

import (
	"fmt"
	"os"
	"testing"
)

func ExampleComponent() {
	log := Component("tag1", "tag2")
	log.Debug("hello world")
	log.Info("hello %v", "world")
	log.Warn("hello", "world")
	log.WithStack(5).Error("%v %v", "hello", "world")

	log.WithTag("tempTag").Debug("I have temp tag")
	log.Warn("hello", "world")

	log.AddTag("tag3")
	log.Debug("hello world")

	// output:
}

func ExampleSetDefaultWriter() {
	SetDefaultWriter(func(o *JsonLogObject) {
		_, _ = fmt.Fprintf(os.Stderr, o.String())
	})

	log := Component()
	log.Debug("hello world")

	// output:
}

func BenchmarkLog(b *testing.B) {
	SetDefaultWriter(func(o *JsonLogObject) {

	})
	log := Component()

	for i := 0; i < b.N; i++ {
		log.Debug("hello world")
	}
}
