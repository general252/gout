package ulog2

func ExampleDebug() {
	Debug("hello %v", "world")

	// output:
}

func ExampleDebugT() {
	tag := Tags("tag1", "tag2")
	DebugT(tag, "hello %v", "world")

	// output:
}
