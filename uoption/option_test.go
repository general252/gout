package uoption

import "log"

type dialOption struct {
	readBufferSize  int
	writeBufferSize int
}

func WithReadBufferSize(s int) Option[*dialOption] {
	return NewFuncOption[*dialOption](func(o *dialOption) {
		o.readBufferSize = s
	})
}

func WithWriteBufferSize(s int) Option[*dialOption] {
	return NewFuncOption[*dialOption](func(o *dialOption) {
		o.writeBufferSize = s
	})
}

func ExampleNewFuncOption() {
	var opts []Option[*dialOption]

	opts = append(opts, WithReadBufferSize(256))
	opts = append(opts, WithWriteBufferSize(512))

	var conf dialOption
	for _, opt := range opts {
		opt.Apply(&conf)
	}

	log.Println(conf)

	// output:

}
