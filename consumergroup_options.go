package historitor

type consumerGroupOptions struct {
	Name    string
	StartAt EntryID
	Members map[string]Consumer
}

func newDefaultConsumerGroupOptions() consumerGroupOptions {
	return consumerGroupOptions{
		StartAt: StartFromBeginning,
		Members: make(map[string]Consumer),
	}
}

var globalConsumerGroupOptions []ConsumerGroupOption

// ConsumerGroupOption is an option for configuring a ConsumerGroup.
type ConsumerGroupOption interface {
	apply(*consumerGroupOptions)
}

// funcConsumerGroupOption is a ConsumerGroupOption that calls a function.
// It is used to wrap a function, so it satisfies the ConsumerGroupOption interface.
type funcConsumerGroupOption struct {
	f func(*consumerGroupOptions)
}

func (fdo *funcConsumerGroupOption) apply(opts *consumerGroupOptions) {
	fdo.f(opts)
}

func newFuncConsumerGroupOption(f func(*consumerGroupOptions)) *funcConsumerGroupOption {
	return &funcConsumerGroupOption{
		f: f,
	}
}

// WithConsumerGroupName returns a ConsumerGroupOption that uses the provided name.
func WithConsumerGroupName(name string) ConsumerGroupOption {
	return newFuncConsumerGroupOption(func(opts *consumerGroupOptions) {
		opts.Name = name
	})
}

// WithConsumerGroupStartAt returns a ConsumerGroupOption that uses the provided start at entry ID.
func WithConsumerGroupStartAt(startAt EntryID) ConsumerGroupOption {
	return newFuncConsumerGroupOption(func(opts *consumerGroupOptions) {
		opts.StartAt = startAt
	})
}

// WithConsumerGroupMember returns a ConsumerGroupOption that uses the provided member.
func WithConsumerGroupMember(member Consumer) ConsumerGroupOption {
	return newFuncConsumerGroupOption(func(opts *consumerGroupOptions) {
		opts.Members[member.GetName()] = member
	})
}
