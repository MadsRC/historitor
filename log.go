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
	ErrNoSuchEntry    = fmt.Errorf("no such entry")
)

type Log struct {
	name             string
	groups           map[string]*consumerGroup
	treeMux          sync.RWMutex
	entries          art.Tree
	firstEntry       EntryID
	lastEntry        EntryID
	maxPendingAge    time.Duration
	maxDeliveryCount int
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
		name:          opts.Name,
		maxPendingAge: opts.MaxPendingAge,
		groups:        make(map[string]*consumerGroup),
		treeMux:       sync.RWMutex{},
		entries:       art.New(),
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
// Once a member reads an Entry, it is added to the Pending Entries List for the consumer group and only removed when the
// member acknowledges the Entry. Entries that are pending will not be returned to any other group member.
//
// If the consumer has pending entries older than [WithMaxPendingAge], up to maxMessages will be returned from the
// pending entries list before reading from the log.
//
// Read is safe for concurrent use.
func (l *Log) Read(g, c string, maxMessages int) ([]Entry, error) {
	group, ok := l.getGroup(g)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchGroup, g)
	}
	consumer, ok := group.getMember(c)
	if !ok {
		return nil, fmt.Errorf("%w in group: %s (group): %s", ErrNoSuchConsumer, g, c)
	}

	out := make([]Entry, 0, maxMessages)

	l.treeMux.RLock()
	defer l.treeMux.RUnlock()

	// check for pending entries
	out, err := l.addPendingEntries(group, *consumer, maxMessages, out)
	if err != nil {
		return nil, err
	}
	if maxMessages > 0 && len(out) >= maxMessages {
		return out, nil
	}
	// no more pending entries, read from log
	out, err = l.addEntries(group, *consumer, maxMessages, out)
	if err != nil {
		return nil, err
	}
	// update the startAt for the group
	group.SetStartAt(out[len(out)-1].ID)

	return out, nil
}

func (l *Log) addPendingEntries(group *consumerGroup, consumer consumerGroupMember, maxMessages int, entries []Entry) ([]Entry, error) {
	for _, pe := range group.getPendingEntriesForConsumer(consumer.name) {
		if time.Since(pe.deliveredAt) > l.maxPendingAge {
			group.incrementDeliveryCountAndTime(pe.id)
			p, ok := l.entries.Search(art.Key(pe.id.String()))
			if !ok {
				return entries, fmt.Errorf("couldn't locate PEL entry in log: %w: %s", ErrNoSuchEntry, pe.id)
			}
			entries = append(entries, Entry{
				ID:      pe.id,
				Payload: p,
			})
			if maxMessages > 0 && len(entries) >= maxMessages {
				break
			}
		}
	}

	return entries, nil
}

func (l *Log) addEntries(group *consumerGroup, consumer consumerGroupMember, maxMessages int, entries []Entry) ([]Entry, error) {
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
			return entries, err
		}

		// check if entry is pending
		_, ok := group.getPendingEntry(eid)
		if ok {
			continue
		}

		// add entry to Pending Entries List
		group.addPendingEntry(eid, consumer.name)
		entries = append(entries, Entry{
			ID:      eid,
			Payload: n.Value(),
		})

		if maxMessages > 0 && len(entries) >= maxMessages {
			break
		}
	}

	return entries, nil
}

func (l *Log) getGroup(name string) (*consumerGroup, bool) {
	l.treeMux.RLock()
	g, ok := l.groups[name]
	l.treeMux.RUnlock()
	return g, ok
}

// Acknowledge acknowledges that a consumer group member has read a log entry. The log entry is removed from the
// consumer group's Pending Entries List.
//
// Acknowledge is safe for concurrent use.
func (l *Log) Acknowledge(g, c string, id EntryID) error {
	group, ok := l.getGroup(g)
	if !ok {
		return fmt.Errorf("%w: %s", ErrNoSuchGroup, g)
	}

	pe, ok := group.getPendingEntry(id)
	if !ok {
		return fmt.Errorf("entry %s not pending", id)
	}
	if pe.consumer != c {
		return fmt.Errorf("%w: entry %s not pending for consumer %s", ErrNoSuchConsumer, id, c)
	}

	group.removePendingEntry(id)

	return nil
}

// Cleanup removes pending entries that are older than [WithMaxPendingAge] and have been delivered more than
// [WithMaxDeliveryCount] times. This effectively releases the log entry back to the consumer group for reading.
//
// Cleanup is safe for concurrent use.
func (l *Log) Cleanup() {
	l.treeMux.Lock()
	defer l.treeMux.Unlock()

	for _, group := range l.groups {
		for _, pe := range group.pel {
			if time.Since(pe.deliveredAt) > l.maxPendingAge && pe.deliveryCount > l.maxDeliveryCount {
				group.removePendingEntry(pe.id)
			}
		}
	}

}
