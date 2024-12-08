package historitor

import (
	art "github.com/plar/go-adaptive-radix-tree"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestEntryID_IsZero(t *testing.T) {
	e := EntryID{}
	require.True(t, e.IsZero())
	require.True(t, ZeroEntryID.IsZero())
}

func TestEntryID(t *testing.T) {
	e := EntryID{
		Time: time.Now().Truncate(time.Millisecond),
		Seq:  123,
	}
	s := e.String()
	e2, err := NewEntryID(s)
	require.NoError(t, err)
	require.Equal(t, e, e2)
}

func TestLog_Read_from_beginning(t *testing.T) {
	l, err := NewLog(WithName("test"))
	require.NoError(t, err)
	tree := art.New()
	keyOne := "1526919030474-55"
	keyTwo := "1526919030474-56"
	keyThree := "1526919030474-57"
	tree.Insert(art.Key(keyOne), "one")
	tree.Insert(art.Key(keyTwo), "two")
	tree.Insert(art.Key(keyThree), "three")
	l.entries = tree

	groupMembers := map[string]consumerGroupMember{
		"consumer1": {
			name: "consumer1",
		},
	}
	groups := map[string]*consumerGroup{
		"group1": {
			name:    "group1",
			members: groupMembers,
			mut:     sync.RWMutex{},
			pel:     make(pendingEntriesList),
			startAt: StartFromBeginning,
		},
	}
	l.groups = groups

	entries, err := l.Read("group1", "consumer1", 3)
	require.NoError(t, err)
	require.Len(t, entries, 3)
}
