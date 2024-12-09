package historitor

import (
	art "github.com/plar/go-adaptive-radix-tree"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestART_insert tests that iterating over the tree after inserting a few keys returns the keys in the correct order.
func TestART_iteration_order(t *testing.T) {
	l := art.New()

	l.Insert([]byte("1234-1"), "value1")
	l.Insert([]byte("1234-2"), "value2")
	l.Insert([]byte("1234-10"), "value3")
	l.Insert([]byte("1234-100"), "value4")

	out := make([]string, 0)
	iter := l.Iterator()
	for iter.HasNext() {
		k, _ := iter.Next()
		out = append(out, string(k.Key()))
	}

	expectedOrder := []string{"1234-1", "1234-2", "1234-10", "1234-100"}

	require.Len(t, out, 4)
	require.Equal(t, expectedOrder, out)

}
