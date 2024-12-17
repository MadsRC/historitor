//go:build !integration

package historitor

import (
	"fmt"
	art "github.com/plar/go-adaptive-radix-tree/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

type stubIterator struct {
	content  []art.Node
	position int
}

func (i *stubIterator) HasNext() bool {
	return i.position < len(i.content)
}

func (i *stubIterator) Next() (art.Node, error) {
	if i.position >= len(i.content) {
		return nil, art.ErrNoMoreNodes
	}
	n := i.content[i.position]
	i.position++
	return n, nil
}

type stubNode struct {
	kind  art.Kind
	key   art.Key
	value art.Value
}

func (n *stubNode) Kind() art.Kind {
	return n.kind
}

func (n *stubNode) Key() art.Key {
	return n.key
}

func (n *stubNode) Value() art.Value {
	return n.value
}

func TestIterateFrom_HasNext(t *testing.T) {
	iter := &stubIterator{
		content: []art.Node{
			&stubNode{},
		},
		position: 0,
	}
	i := newIterateFrom(nil, iter)
	require.True(t, i.HasNext())
}

// Make sure that the key designated using the iterateFrom's key attribute is not included in the iteration.
// In this case, the tree only has one node, and the key is the same as the key we are iterating from, so
// when we call Next, we should get [ErrNoMoreEntries].
func TestIterateFrom_Next_key_not_inclusive(t *testing.T) {
	e1 := stubNode{kind: art.Leaf, key: []byte("key"), value: []byte("value")}
	iter := &stubIterator{content: []art.Node{&e1}}
	i := newIterateFrom([]byte("key"), iter)
	out, err := i.Next()
	require.Nil(t, out)
	require.ErrorIs(t, err, ErrNoMoreEntries)
}

// Make sure that the key designated using the iterateFrom's key attribute is not included in the iteration.
// In this case, the tree has two nodes, and the second node's key is the same as the key we are iterating from, so
// when we call Next, we should get [ErrNoMoreEntries].
func TestIterateFrom_Next_key_not_inclusive2(t *testing.T) {
	e1 := stubNode{kind: art.Leaf, key: []byte("key"), value: []byte("value")}
	e2 := stubNode{kind: art.Leaf, key: []byte("key2"), value: []byte("value2")}
	iter := &stubIterator{content: []art.Node{&e1, &e2}}
	i := newIterateFrom([]byte("key2"), iter)
	out, err := i.Next()
	require.Nil(t, out)
	require.ErrorIs(t, err, ErrNoMoreEntries)
}

// Make sure that calling Next on an empty tree returns [ErrNoMoreEntries].
func TestIterateFrom_Next_empty(t *testing.T) {
	iter := &stubIterator{}
	i := newIterateFrom(nil, iter)
	out, err := i.Next()
	require.Nil(t, out)
	require.ErrorIs(t, err, ErrNoMoreEntries)
}

// When calling next on an interator with no key specified, the first node should be returned.
func TestIterateFrom_Next_nil_key(t *testing.T) {
	e1 := stubNode{kind: art.Leaf, key: []byte("key"), value: []byte("value")}
	e2 := stubNode{kind: art.Leaf, key: []byte("key2"), value: []byte("value2")}
	iter := &stubIterator{content: []art.Node{&e1, &e2}}
	i := newIterateFrom(nil, iter)
	out, err := i.Next()
	require.NoError(t, err)
	require.Equal(t, &e1, out)
}

func TestIterateFrom_Next_key_not_found(t *testing.T) {
	e1 := stubNode{kind: art.Leaf, key: []byte("key"), value: []byte("value")}
	e2 := stubNode{kind: art.Leaf, key: []byte("key2"), value: []byte("value2")}
	iter := &stubIterator{content: []art.Node{&e1, &e2}}
	i := newIterateFrom([]byte("key3"), iter)
	out, err := i.Next()
	require.Nil(t, out)
	require.ErrorIs(t, err, ErrNoMoreEntries)
}

// iterate through the tree until the key is found, then continue to return the next node until there are no more nodes.
func TestIterateFrom_Next(t *testing.T) {
	e1 := stubNode{kind: art.Leaf, key: []byte("key"), value: []byte("value")}
	e2 := stubNode{kind: art.Leaf, key: []byte("key2"), value: []byte("value2")}
	iter := &stubIterator{content: []art.Node{&e1, &e2}}
	i := newIterateFrom(nil, iter)
	out, err := i.Next()
	require.NoError(t, err)
	require.Equal(t, &e1, out)
	out, err = i.Next()
	require.NoError(t, err)
	require.Equal(t, &e2, out)
	out, err = i.Next()
	require.Nil(t, out)
	require.ErrorIs(t, err, ErrNoMoreEntries)
}

func TestNextEntryError_convert_ErrNoMoreNodes_to_ErrNoMoreEntries(t *testing.T) {
	err := nextEntryError(art.ErrNoMoreNodes, "test")
	require.ErrorIs(t, err, ErrNoMoreEntries)
}

func TestNextEntryError(t *testing.T) {
	err := nextEntryError(ErrNoSuchConsumer, "test")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf("test: %s", ErrNoSuchConsumer), err.Error())

	t.Run("no wrap", func(t *testing.T) {
		require.NotErrorIs(t, err, ErrNoSuchConsumer)
	})
}
