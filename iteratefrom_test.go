package historitor

import (
	mock_art "github.com/MadsRC/historitor/mocks/github.com/plar/go-adaptive-radix-tree"
	art "github.com/plar/go-adaptive-radix-tree"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testKey = art.Key("key")
)

func TestIterateFrom_HasNext(t *testing.T) {
	mockIter := mock_art.NewMockIterator(t)
	mockIter.EXPECT().HasNext().Return(true)

	iter := newIterateFrom(testKey, mockIter)
	require.True(t, iter.HasNext())
}

// Make sure that the key designated using the iterateFrom's key attribute is not included in the iteration.
// In thise case, the tree only has one node, and the key is the same as the key we are iterating from, so
// when we call Next, we should get an error.
func TestIterateFrom_Next_key_not_inclusive(t *testing.T) {
	mockIter := mock_art.NewMockIterator(t)
	mockNode := mock_art.NewMockNode(t)
	mockIter.EXPECT().Next().Return(mockNode, nil)
	mockNode.EXPECT().Key().Return(testKey)
	mockIter.EXPECT().HasNext().Return(false)

	iter := newIterateFrom(testKey, mockIter)
	n, err := iter.Next()
	require.Nil(t, n)
	require.ErrorIs(t, err, ErrNoMoreEntries)
}

// iterate through the tree until the key is found, then continue to return the next node until there are no more nodes.
func TestIterateFrom_Next(t *testing.T) {
	var keys = []art.Key{art.Key("key1"), art.Key("key2"), art.Key("key3"), art.Key("key4")}
	mockIter := mock_art.NewMockIterator(t)
	mockNodeOne := mock_art.NewMockNode(t)
	mockNodeTwo := mock_art.NewMockNode(t)
	mockNodeThree := mock_art.NewMockNode(t)
	mockNodeFour := mock_art.NewMockNode(t)

	mockIter.EXPECT().Next().Return(mockNodeOne, nil).Once()
	mockNodeOne.EXPECT().Key().Return(keys[0]).Once()
	mockIter.EXPECT().HasNext().Return(true).Once()

	mockIter.EXPECT().Next().Return(mockNodeTwo, nil).Once()
	mockNodeTwo.EXPECT().Key().Return(keys[1]).Once()
	mockIter.EXPECT().HasNext().Return(true).Once()

	mockIter.EXPECT().Next().Return(mockNodeThree, nil).Once()

	iter := newIterateFrom(keys[1], mockIter)
	n, err := iter.Next()
	require.NoError(t, err)
	require.Equal(t, mockNodeThree, n)

	mockIter.EXPECT().Next().Return(mockNodeFour, nil).Once()

	n, err = iter.Next()
	require.NoError(t, err)
	require.Equal(t, mockNodeFour, n)

	mockIter.EXPECT().Next().Return(nil, art.ErrNoMoreNodes).Once()
	n, err = iter.Next()
	require.Nil(t, n)
	require.ErrorIs(t, err, ErrNoMoreEntries)
}
