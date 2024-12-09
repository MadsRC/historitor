package historitor

type consumerOptions struct {
	Name string
}

var defaultConsumerOptions = consumerOptions{}

var globalConsumerOptions []ConsumerOption

// ConsumerOption is an option for configuring a Consumer.
type ConsumerOption interface {
	apply(*consumerOptions)
}

// funcConsumerOption is a ConsumerOption that calls a function.
// It is used to wrap a function, so it satisfies the ConsumerOption interface.
type funcConsumerOption struct {
	f func(*consumerOptions)
}

func (fdo *funcConsumerOption) apply(opts *consumerOptions) {
	fdo.f(opts)
}

func newFuncConsumerOption(f func(*consumerOptions)) *funcConsumerOption {
	return &funcConsumerOption{
		f: f,
	}
}

// WithConsumerName returns a ConsumerOption that uses the provided name.
func WithConsumerName(name string) ConsumerOption {
	return newFuncConsumerOption(func(opts *consumerOptions) {
		opts.Name = name
	})
}
