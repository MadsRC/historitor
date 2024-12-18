package historitor

import (
	"time"
)

type logOptions struct {
	// Name is the name of the log.
	Name                   string
	MaxPendingAge          time.Duration
	MaxDeliveryCount       int
	AttemptRedeliveryAfter time.Duration
}

var defaultLogOptions = logOptions{
	MaxPendingAge:          4 * time.Second,
	MaxDeliveryCount:       3,
	AttemptRedeliveryAfter: time.Second,
}

var GlobalLogOptions []LogOption

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

// WithLogName sets the name of the log to the provided name.
func WithLogName(name string) LogOption {
	return newFuncLogOption(func(opts *logOptions) {
		opts.Name = name
	})
}

// WithLogMaxPendingAge sets the maximum age of a log entry before it is considered stale and should be removed from
// the Pending Entries List. This will allow other consumers in the group to attempt to process the log entry.
func WithLogMaxPendingAge(maxPendingAge time.Duration) LogOption {
	return newFuncLogOption(func(opts *logOptions) {
		opts.MaxPendingAge = maxPendingAge
	})
}

// WithLogMaxDeliveryCount sets the maximum number of times re-delivery of a log entry is attempted.
func WithLogMaxDeliveryCount(maxDeliveryCount int) LogOption {
	return newFuncLogOption(func(opts *logOptions) {
		opts.MaxDeliveryCount = maxDeliveryCount
	})
}

// WithLogAttemptRedeliveryAfter sets the duration after which a log entry should be re-delivered to the Consumer if it
// has not been acknowledged.
func WithLogAttemptRedeliveryAfter(attemptRedeliveryAfter time.Duration) LogOption {
	return newFuncLogOption(func(opts *logOptions) {
		opts.AttemptRedeliveryAfter = attemptRedeliveryAfter
	})
}
