//go:build !integration

package historitor

import (
	art "github.com/plar/go-adaptive-radix-tree"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestART_insert tests that iterating over the tree after inserting a few keys returns the keys in the correct order.
func TestART_iteration_order(t *testing.T) {
	l := art.New()

	ti := time.Now().Truncate(time.Millisecond)
	idOne := EntryID{time: ti, seq: 0}
	idTwo := EntryID{time: ti, seq: 1}
	idThree := EntryID{time: ti, seq: 2}
	idFour := EntryID{time: ti, seq: 3}

	l.Insert([]byte(idOne.String()), "value1")
	l.Insert([]byte(idTwo.String()), "value2")
	l.Insert([]byte(idThree.String()), "value3")
	l.Insert([]byte(idFour.String()), "value4")

	out := make([]string, 0)
	iter := l.Iterator()
	for iter.HasNext() {
		k, _ := iter.Next()
		out = append(out, string(k.Key()))
	}

	expectedOrder := []string{
		idOne.String(),
		idTwo.String(),
		idThree.String(),
		idFour.String(),
	}

	require.Len(t, out, 4)
	require.Equal(t, expectedOrder, out)

}
