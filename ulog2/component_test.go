package ulog2

import (
	"fmt"
)

func ExampleComponent() {
	log := Component("tag1", "tag2")
	log.Debug("hello world")
	log.Info("hello %v", "world")
	log.Warn("hello", "world")
	log.Error("%v %v", "hello", "world")

	log.WithTag("tempTag").Debug("I have temp tag")

	log.AddTag("tag3")
	log.Debug("hello world")

	// output:
	// 2023-03-23 20:22:33 component_test.go:5 [D] [tag1, tag2] hello world
	// 2023-03-23 20:22:33 component_test.go:6 [I] [tag1, tag2] hello world
	// 2023-03-23 20:22:33 component_test.go:7 [W] [tag1, tag2] hello world
	// 2023-03-23 20:22:33 component_test.go:8 [E] [tag1, tag2] hello world
	// 2023-03-23 20:22:33 component_test.go:10 [D] [tag1, tag2, tempTag] I have temp tag
	// 2023-03-23 20:22:33 component_test.go:13 [D] [tag1, tag2, tag3] hello world
}

func ExampleSetDefaultWriter() {
	SetDefaultWriter(func(o *JsonLogObject) {
		fmt.Println(o.String())
	})

	log := Component()
	log.Debug("hello world")

	// output:
}
