package server

type option struct {
	convert
	object
}

var options = &option{
	convert: &rawConvert{},
	object:  &ocppType{},
}

type opt func(*option)

func WithOptions(opts ...opt) {
	for _, opt := range opts {
		opt(options)
	}
}
