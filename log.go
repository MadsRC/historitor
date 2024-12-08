package historitor

import (
	"fmt"
	"github.com/plar/go-adaptive-radix-tree"
	"sync"
	"time"
)

var (
	ErrNoSuchGroup    = fmt.Errorf("no such consumer group")
	ErrNoSuchConsumer = fmt.Errorf("no such consumer")
)

type Log struct {
	name       string
	groups     map[string]consumerGroup
	treeMux    sync.RWMutex
	entries    art.Tree
	firstEntry EntryID
	lastEntry  EntryID
}

func NewLog(options ...LogOption) (*Log, error) {
	opts := defaultLogOptions
	for _, opt := range globalLogOptions {
		opt.apply(&opts)
	}
	for _, opt := range options {
		opt.apply(&opts)
	}

	return &Log{
		name:    opts.Name,
		groups:  make(map[string]consumerGroup),
		treeMux: sync.RWMutex{},
		entries: art.New(),
	}, nil
}

// Write writes a new log entry to the log. It returns the ID of the log entry.
//
// Write is safe for concurrent use.
func (l *Log) Write(payload any) string {
	l.treeMux.Lock()
	id := EntryID{
		Time: time.Now(),
	}
	l.write(&id, payload)
	l.treeMux.Unlock()
	return id.String()
}

// write is not safe for concurrent use. It should be called with the treeMux locked.
// write is a recursive function that will attempt to write a log entry to the log. If the key already exists, it will
// increment the sequence number and try again by calling itself.
func (l *Log) write(id *EntryID, payload any) {
	ov, upd := l.entries.Insert(art.Key(id.String()), payload)
	if upd {
		// restore the value we just overwrote
		l.entries.Insert(art.Key(id.String()), ov)
		// increment the sequence number and try again
		id.Seq++
		l.write(id, payload)
	}
	l.lastEntry = *id
	return
}

// Read reads up to maxMessages log entries from the log. If maxMessages is 0, it will read all log entries.
// Returning an empty slice means there are no log entries to read.
// Group and consumer name are used to track which log entries have been read by which consumer group members.
// If a consumer group member has read a log entry, it will not be returned to any other group member.
// Once a member reads an event, it is added to the Pending Events List for the consumer group and only removed when the
// member acknowledges the event. Events that are pending will not be returned to any other group member.
//
// Read is safe for concurrent use.
func (l *Log) Read(g, c string, maxMessages int) ([]Entry, error) {
	group, ok := l.getGroup(g)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchGroup, g)
	}
	_, ok = group.getMember(c)
	if !ok {
		return nil, fmt.Errorf("%w in group: %s (group): %s", ErrNoSuchConsumer, g, c)
	}

	out := make([]Entry, 0, maxMessages)

	l.treeMux.RLock()
	defer l.treeMux.RUnlock()
	var iter art.Iterator
	if group.GetStartAt() == StartFromBeginning {
		iter = l.entries.Iterator()
	} else {
		iter = newIterateFrom(art.Key(group.GetStartAt().String()), l.entries.Iterator())
	}
	for iter.HasNext() {
		n, err := iter.Next()
		if err != nil {
			break
		}

		eid, err := NewEntryID(string(n.Key()))
		if err != nil {
			return nil, err
		}

		// check if entry is pending
		_, ok := group.getPendingEvent(eid)
		if ok {
			continue
		}

		// add entry to pending events list
		group.addPendingEvent(eid, c)
		out = append(out, Entry{
			ID:      eid,
			Payload: n.Value(),
		})

		if maxMessages > 0 && len(out) >= maxMessages {
			break
		}
	}

	// update the startAt for the group
	group.SetStartAt(out[len(out)-1].ID)

	return out, nil
}

func (l *Log) getGroup(name string) (*consumerGroup, bool) {
	l.treeMux.RLock()
	g, ok := l.groups[name]
	l.treeMux.RUnlock()
	return &g, ok
}
