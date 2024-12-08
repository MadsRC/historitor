package historitor

import "time"

type logOptions struct {
	// Name is the name of the log.
	Name             string
	MaxPendingAge    time.Duration
	MaxDeliveryCount int
}

var defaultLogOptions = logOptions{
	MaxPendingAge: time.Second,
}

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

// WithMaxPendingAge returns a LogOption that uses the provided max pending age.
func WithMaxPendingAge(maxPendingAge time.Duration) LogOption {
	return newFuncLogOption(func(opts *logOptions) {
		opts.MaxPendingAge = maxPendingAge
	})
}

// WithMaxDeliveryCount returns a LogOption that uses the provided max delivery count.
func WithMaxDeliveryCount(maxDeliveryCount int) LogOption {
	return newFuncLogOption(func(opts *logOptions) {
		opts.MaxDeliveryCount = maxDeliveryCount
	})
}
