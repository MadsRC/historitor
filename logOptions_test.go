package historitor

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFuncLogOption_apply(t *testing.T) {
	flo := funcLogOption{
		f: func(opts *logOptions) {
			opts.Name = "test"
		},
	}
	opts := logOptions{}
	flo.apply(&opts)
	require.Equal(t, "test", opts.Name)
}

func TestNewFuncLogOption(t *testing.T) {
	f := func(opts *logOptions) {
		opts.Name = "test"
	}
	flo := newFuncLogOption(f)
	require.NotNil(t, flo.f)
}

func TestWithLogName(t *testing.T) {
	opts := logOptions{}
	lo := WithLogName("test")
	lo.apply(&opts)
	require.Equal(t, "test", opts.Name)
}

func TestWithLogMaxPendingAge(t *testing.T) {
	opts := logOptions{}
	lo := WithLogMaxPendingAge(4)
	lo.apply(&opts)
	require.Equal(t, time.Duration(4), opts.MaxPendingAge)
}

func TestWithLogMaxDeliveryCount(t *testing.T) {
	opts := logOptions{}
	lo := WithLogMaxDeliveryCount(3)
	lo.apply(&opts)
	require.Equal(t, 3, opts.MaxDeliveryCount)
}

func TestWithLogAttemptRedeliveryAfter(t *testing.T) {
	opts := logOptions{}
	lo := WithLogAttemptRedeliveryAfter(1)
	lo.apply(&opts)
	require.Equal(t, time.Duration(1), opts.AttemptRedeliveryAfter)
}
