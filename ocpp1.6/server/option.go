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

//WithOptions used for some custom events
func WithOptions(opts ...opt) {
	for _, opt := range opts {
		opt(options)
	}
}
