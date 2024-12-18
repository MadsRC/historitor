//go:build !integration

package historitor

import (
	art "github.com/plar/go-adaptive-radix-tree/v2"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestNewLog(t *testing.T) {
	GlobalLogOptions = append(GlobalLogOptions, WithLogMaxPendingAge(1))
	defer func() {
		GlobalLogOptions = GlobalLogOptions[:len(GlobalLogOptions)-1]
	}()
	l, err := NewLog(WithLogName("test"))
	require.NoError(t, err)
	require.NotNil(t, l)
	require.Equal(t, "test", l.name)
	require.Equal(t, time.Duration(1), l.maxPendingAge)
	require.NotNil(t, l.entries)
	require.NotNil(t, l.groups)
}

func TestLog_Size(t *testing.T) {
	l := &Log{
		entries: art.New(),
	}
	require.Equal(t, 0, l.Size())
	l.entries.Insert(art.Key("key"), "value")
	require.Equal(t, 1, l.Size())
}

func TestLog_Write_id_has_timezone_set_to_utc(t *testing.T) {
	l, err := NewLog(WithLogName(t.Name()))
	require.NoError(t, err)

	id := l.Write("value")
	require.Equal(t, "UTC", id.time.Location().String())
}

func TestLog_write_key_already_exists(t *testing.T) {
	id := fakeTestEntryID1
	l := &Log{
		entries: art.New(),
	}
	l.entries.Insert(art.Key(id.String()), "value")
	l.write(&id, "value")
	require.Equal(t, 2, l.entries.Size())
	e1, ok := l.entries.Search(art.Key(fakeTestEntryID1.String()))
	require.True(t, ok)
	require.Equal(t, "value", e1)
	require.Equal(t, uint64(1), fakeTestEntryID1.seq)
	e2, ok := l.entries.Search(art.Key(id.String()))
	require.True(t, ok)
	require.Equal(t, "value", e2)
	require.Equal(t, uint64(2), id.seq)
}

func TestLog_Read_id_has_timezone_set_to_utc(t *testing.T) {
	tree := art.New()
	keyOne := fakeTestEntryID1.String()
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

func TestLog_Read_from_beginning(t *testing.T) {
	l, err := NewLog(WithLogName("test"))
	require.NoError(t, err)
	tree := art.New()
	keyOne := fakeTestEntryID1.String()
	keyTwo := fakeTestEntryID2.String()
	keyThree := fakeTestEntryID3.String()
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
