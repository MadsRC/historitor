package historitor

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/plar/go-adaptive-radix-tree"
	"sync"
	"time"
)

var (
	ErrNoSuchGroup    = fmt.Errorf("no such Consumer group")
	ErrNoSuchConsumer = fmt.Errorf("no such Consumer")
	ErrNoSuchEntry    = fmt.Errorf("no such entry")
)

// externalLog is used to represent a Log in a way that can easily be encoded and decoded using the gob package.
type externalLog struct {
	Name                   string
	Groups                 map[string]*ConsumerGroup
	Entries                []Entry
	FirstEntry             EntryID
	LastEntry              EntryID
	MaxPendingAge          time.Duration
	MaxDeliveryCount       int
	AttemptRedeliveryAfter time.Duration
}

// Log is a transactional log that allows for multiple readers and writers. It is backed by an in-memory radix tree.
//
// Instances of [Log] should have their [Log.Cleanup] method called periodically to ensure that non-acknowledged log
// entries are released for re-delivery.
type Log struct {
	name                   string
	groups                 map[string]*ConsumerGroup
	treeMux                sync.RWMutex
	entries                art.Tree
	firstEntry             EntryID
	lastEntry              EntryID
	maxPendingAge          time.Duration
	maxDeliveryCount       int
	attemptRedeliveryAfter time.Duration
}

// NewLog creates a new log with the provided options.
func NewLog(options ...LogOption) (*Log, error) {
	opts := defaultLogOptions
	for _, opt := range globalLogOptions {
		opt.apply(&opts)
	}
	for _, opt := range options {
		opt.apply(&opts)
	}

	return &Log{
		name:                   opts.Name,
		maxPendingAge:          opts.MaxPendingAge,
		maxDeliveryCount:       opts.MaxDeliveryCount,
		attemptRedeliveryAfter: opts.AttemptRedeliveryAfter,
		groups:                 make(map[string]*ConsumerGroup),
		treeMux:                sync.RWMutex{},
		entries:                art.New(),
	}, nil
}

// Size returns the number of log entries in the log.
func (l *Log) Size() int {
	return l.entries.Size()
}

// Write writes a new log entry to the log. It returns the ID of the log entry.
//
// Write is safe for concurrent use.
func (l *Log) Write(payload any) EntryID {
	l.treeMux.Lock()
	id := NewEntryID(time.Now().Truncate(time.Millisecond).UTC(), 0)
	l.write(&id, payload)
	l.treeMux.Unlock()
	return id
}

// write is not safe for concurrent use. It should be called with the treeMux locked.
// write is a recursive function that will attempt to write a log entry to the log. If the key already exists, it will
// increment the sequence number and try again by calling itself.
//
// The time of the EntryID is truncated to milliseconds.
func (l *Log) write(id *EntryID, payload any) {
	id.time = id.time.Truncate(time.Millisecond)
	ov, upd := l.entries.Insert(art.Key(id.String()), payload)
	if upd {
		// restore the value we just overwrote
		l.entries.Insert(art.Key(id.String()), ov)
		// increment the sequence number and try again
		id.seq++
		l.write(id, payload)
	}
	l.lastEntry = *id
	return
}

// Read reads up to maxMessages log entries from the log. If maxMessages is 0, it will read all log entries.
// Returning an empty slice means there are no log entries to read.
// Group and Consumer name are used to track which log entries have been read by which Consumer group members.
// If a Consumer group member has read a log entry, it will not be returned to any other group member.
// Once a member reads an Entry, it is added to the Pending Entries List for the Consumer group and only removed when the
// member acknowledges the Entry. Entries that are pending will not be returned to any other group member.
//
// If the Consumer has pending entries older than [WithLogAttemptRedeliveryAfter] and that have been delivered more than
// [WithLogMaxDeliveryCount], up to maxMessages will be returned from the pending entries list before reading from the
// log.
//
// If there are no more events to read from the log, the method will return an empty slice.
//
// Read is safe for concurrent use.
func (l *Log) Read(g, c string, maxMessages int) ([]Entry, error) {
	group, ok := l.getGroup(g)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchGroup, g)
	}
	consumer, ok := group.GetMember(c)
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
		if !errors.Is(err, art.ErrNoMoreNodes) {
			return nil, err
		}
	}
	if len(out) > 0 && !errors.Is(err, art.ErrNoMoreNodes) {
		// update the startAt for the group to the last entry read if we actually read something off the log
		group.SetStartAt(out[len(out)-1].ID)
	}

	return out, nil
}

func (l *Log) addPendingEntries(group *ConsumerGroup, consumer Consumer, maxMessages int, entries []Entry) ([]Entry, error) {
	for _, pe := range group.GetPendingEntriesForConsumer(consumer.name) {
		if time.Since(pe.DeliveredAt) > l.attemptRedeliveryAfter && pe.DeliveryCount < l.maxDeliveryCount {
			group.incrementDeliveryCountAndTime(pe.ID)
			p, ok := l.entries.Search(art.Key(pe.ID.String()))
			if !ok {
				return entries, fmt.Errorf("couldn't locate PEL entry in log: %w: %s", ErrNoSuchEntry, pe.ID)
			}
			entries = append(entries, Entry{
				ID:      pe.ID,
				Payload: p,
			})
			if maxMessages > 0 && len(entries) >= maxMessages {
				break
			}
		}
	}

	return entries, nil
}

func (l *Log) addEntries(group *ConsumerGroup, consumer Consumer, maxMessages int, entries []Entry) ([]Entry, error) {
	var iter art.Iterator
	if group.GetStartAt() == StartFromBeginning {
		iter = l.entries.Iterator()
	} else {
		iter = newIterateFrom(art.Key(group.GetStartAt().String()), l.entries.Iterator())
	}
	for iter.HasNext() {
		n, err := iter.Next()
		if err != nil {
			return nil, err
		}

		eid, err := ParseEntryID(string(n.Key()))
		if err != nil {
			return entries, err
		}

		// check if entry is pending
		_, ok := group.GetPendingEntry(eid)
		if ok {
			continue
		}

		// add entry to Pending Entries List
		group.AddPendingEntry(eid, consumer.name)
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

func (l *Log) getGroup(name string) (*ConsumerGroup, bool) {
	l.treeMux.RLock()
	g, ok := l.groups[name]
	l.treeMux.RUnlock()
	return g, ok
}

// AddGroup adds a Consumer group to the log.
func (l *Log) AddGroup(group *ConsumerGroup) {
	l.treeMux.Lock()
	l.groups[group.name] = group
	l.treeMux.Unlock()
}

// RemoveGroup removes a Consumer group from the log.
func (l *Log) RemoveGroup(name string) {
	l.treeMux.Lock()
	delete(l.groups, name)
	l.treeMux.Unlock()
}

// ListGroups returns a list of all Consumer groups.
func (l *Log) ListGroups() []*ConsumerGroup {
	l.treeMux.RLock()
	defer l.treeMux.RUnlock()

	out := make([]*ConsumerGroup, 0, len(l.groups))
	for _, g := range l.groups {
		out = append(out, g)
	}
	return out
}

// Acknowledge acknowledges that a Consumer group member has read a log entry. The log entry is removed from the
// Consumer group's Pending Entries List.
//
// Acknowledge is safe for concurrent use.
func (l *Log) Acknowledge(g, c string, id EntryID) error {
	group, ok := l.getGroup(g)
	if !ok {
		return fmt.Errorf("%w: %s", ErrNoSuchGroup, g)
	}

	pe, ok := group.GetPendingEntry(id)
	if !ok {
		return fmt.Errorf("entry %s not pending", id)
	}
	if pe.Consumer != c {
		return fmt.Errorf("%w: entry %s not pending for Consumer %s", ErrNoSuchConsumer, id, c)
	}

	group.RemovePendingEntry(id)

	return nil
}

// Cleanup runs a series of housekeeping actions on the log.
//
//   - Remove pending entries that are older than [WithLogMaxPendingAge] to allow other consumers to attempt to process
//     the log entry.
//   - Remove pending entries that have been delivered more than [WithLogMaxDeliveryCount] times and are older than
//     [WithLogAttemptRedeliveryAfter].
//
// Cleanup is safe for concurrent use.
func (l *Log) Cleanup() {
	l.treeMux.Lock()
	defer l.treeMux.Unlock()

	for _, group := range l.groups {
		for _, pe := range group.pel {
			if time.Since(pe.DeliveredAt) > l.attemptRedeliveryAfter && pe.DeliveryCount > l.maxDeliveryCount {
				group.RemovePendingEntry(pe.ID)
			} else if time.Since(pe.DeliveredAt) > l.maxPendingAge {
				group.RemovePendingEntry(pe.ID)
			}
		}
	}
}

// UpdateEntry updates the payload of a log entry. If the log entry does not exist, it will return false.
//
// UpdateEntry is safe for concurrent use.
func (l *Log) UpdateEntry(id EntryID, payload any) bool {
	l.treeMux.Lock()
	defer l.treeMux.Unlock()

	_, upd := l.entries.Insert(art.Key(id.String()), payload)
	if !upd {
		l.entries.Delete(art.Key(id.String()))
	}

	return upd
}

// MarshalBinary encodes a Log into a gob-encoded byte slice.
func (l *Log) MarshalBinary() ([]byte, error) {
	l.treeMux.RLock()
	defer l.treeMux.RUnlock()
	el := externalLog{
		Name:                   l.name,
		Groups:                 l.groups,
		Entries:                make([]Entry, 0),
		FirstEntry:             l.firstEntry,
		LastEntry:              l.lastEntry,
		MaxPendingAge:          l.maxPendingAge,
		MaxDeliveryCount:       l.maxDeliveryCount,
		AttemptRedeliveryAfter: l.attemptRedeliveryAfter,
	}
	l.entries.ForEach(func(node art.Node) (cont bool) {
		id, err := ParseEntryID(string(node.Key()))
		if err != nil {
			return false
		}
		el.Entries = append(el.Entries, Entry{
			ID:      id,
			Payload: node.Value(),
		})
		return true
	})
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(el)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary decodes a gob-encoded byte slice into a Log.
func (l *Log) UnmarshalBinary(data []byte) error {
	l.treeMux.Lock()
	defer l.treeMux.Unlock()
	dec := gob.NewDecoder(bytes.NewReader(data))
	var el externalLog
	err := dec.Decode(&el)
	if err != nil {
		return err
	}
	l.name = el.Name
	l.groups = el.Groups
	l.firstEntry = el.FirstEntry
	l.lastEntry = el.LastEntry
	l.maxPendingAge = el.MaxPendingAge
	l.maxDeliveryCount = el.MaxDeliveryCount
	l.attemptRedeliveryAfter = el.AttemptRedeliveryAfter
	l.entries = art.New()
	for _, e := range el.Entries {
		l.entries.Insert([]byte(e.ID.String()), e.Payload)
	}

	return nil
}
