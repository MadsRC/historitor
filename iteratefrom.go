package historitor

import (
	art "github.com/plar/go-adaptive-radix-tree"
)

// Ensure iterateFrom implements art.Iterator at compile time
var _ art.Iterator = (*iterateFrom)(nil)

// iterateFrom is an iterator that iterates over a tree starting from a given key. It implements the art.Iterator.
// interface.
// iterateFrom is not inclusive of the key it starts from.
type iterateFrom struct {
	key            art.Key
	iter           art.Iterator
	keyEncountered bool
}

// newIterateFrom creates a new iterateFrom iterator. It takes a key to start from and an art.Iterator to iterate over.
func newIterateFrom(key art.Key, iter art.Iterator) *iterateFrom {
	return &iterateFrom{
		key:  key,
		iter: iter,
	}
}

// HasNext returns true if there are more nodes to iterate over.
func (i *iterateFrom) HasNext() bool {
	return i.iter.HasNext()
}

// Next returns the next node in the iteration. Next will skip nodes until it finds the node with the key that was
// provided when creating the iterator. If the key is not found, Next will return art.ErrNoMoreNodes.
func (i *iterateFrom) Next() (art.Node, error) {
	n, err := i.iter.Next()
	if err != nil {
		return nil, err
	}
	if i.keyEncountered {
		return n, nil
	}
	if string(n.Key()) == string(i.key) {
		i.keyEncountered = true
		if !i.iter.HasNext() {
			return nil, art.ErrNoMoreNodes
		}
		n, err = i.iter.Next()
		if err != nil {
			return nil, err
		}

		return n, nil
	}
	for i.iter.HasNext() {
		val, err := i.iter.Next()
		if err != nil {
			return nil, err
		}
		if string(val.Key()) == string(i.key) {
			i.keyEncountered = true
			return val, nil
		}
	}

	return nil, art.ErrNoMoreNodes
}
