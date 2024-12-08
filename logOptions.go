package historitor

type logOptions struct {
	Name string
}

var defaultLogOptions = logOptions{}

var globalLogOptions []LogOption

// LogOption is an option for configuring a Log.
type LogOption interface {
	apply(*logOptions)
}

// funcLogOption is a LogOption that calls a function.
// It is used to wrap a function, so it satisfies the LogOption interface.
type funcLogOption struct {
	f func(*logOptions)
}

func (fdo *funcLogOption) apply(opts *logOptions) {
	fdo.f(opts)
}

func newFuncLogOption(f func(*logOptions)) *funcLogOption {
	return &funcLogOption{
		f: f,
	}
}

// WithName returns a LogOption that uses the provided name.
func WithName(name string) LogOption {
	return newFuncLogOption(func(opts *logOptions) {
		opts.Name = name
	})
}
