//go:build !integration

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
		time: time.Now().Truncate(time.Millisecond).UTC(),
		seq:  123,
	}
	s := e.String()
	e2, err := ParseEntryID(s)
	require.NoError(t, err)
	require.Equal(t, e, e2)
}

func TestLog_Read_from_beginning(t *testing.T) {
	l, err := NewLog(WithLogName("test"))
	require.NoError(t, err)
	tree := art.New()
	keyOne := "1526919030474-55"
	keyTwo := "1526919030474-56"
	keyThree := "1526919030474-57"
	tree.Insert(art.Key(keyOne), "one")
	tree.Insert(art.Key(keyTwo), "two")
	tree.Insert(art.Key(keyThree), "three")
	l.entries = tree

	groupMembers := map[string]Consumer{
		"consumer1": {
			name: "consumer1",
		},
	}
	groups := map[string]*ConsumerGroup{
		"group1": {
			name:    "group1",
			members: groupMembers,
			mut:     sync.RWMutex{},
			pel:     make(PendingEntriesList),
			startAt: StartFromBeginning,
		},
	}
	l.groups = groups

	entries, err := l.Read("group1", "consumer1", 3)
	require.NoError(t, err)
	require.Len(t, entries, 3)
}

func TestLog_Write_id_has_timezone_set_to_utc(t *testing.T) {
	l, err := NewLog(WithLogName(t.Name()))
	require.NoError(t, err)

	id := l.Write("value")
	require.Equal(t, "UTC", id.time.Location().String())
}

func TestLog_Read_id_has_timezone_set_to_utc(t *testing.T) {
	tree := art.New()
	keyOne := "1526919030474-55"
	tree.Insert(art.Key(keyOne), "value")
	c := NewConsumer(WithConsumerName(t.Name()))
	cg := NewConsumerGroup(WithConsumerGroupName(t.Name()), WithConsumerGroupMember(c))
	l, err := NewLog(WithLogName(t.Name()))
	require.NoError(t, err)
	l.AddGroup(cg)

	l.entries = tree

	entries, err := l.Read(cg.GetName(), c.GetName(), 1)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, "UTC", entries[0].ID.time.Location().String())
}
