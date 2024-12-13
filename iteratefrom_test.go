package historitor

import (
	art "github.com/plar/go-adaptive-radix-tree"
	"github.com/stretchr/testify/require"
	"testing"
)

type fakeIter struct {
	HasNextFn     func() bool
	HasNextCalled bool
	NextFn        func() (art.Node, error)
	NextCalled    bool
}

func (f *fakeIter) HasNext() bool {
	f.HasNextCalled = true
	return f.HasNextFn()
}

func (f *fakeIter) Next() (art.Node, error) {
	f.NextCalled = true
	return f.NextFn()
}

func TestIterateFrom_HasNext(t *testing.T) {
	mock := &fakeIter{
		HasNextFn: func() bool {
			return true
		},
	}

	iter := newIterateFrom(art.Key("key"), mock)
	require.True(t, iter.HasNext())
	require.True(t, mock.HasNextCalled)
}
