package historitor

import (
	"sync"
	"time"
)

// ConsumerGroup is a group of consumers that consume log entries together.
//
// ConsumerGroup must not be copied.
type ConsumerGroup struct {
	name    string
	members map[string]Consumer
	mut     sync.RWMutex
	pel     PendingEntriesList
	startAt EntryID
}

// NewConsumerGroup creates a new Consumer group with the provided options.
func NewConsumerGroup(options ...ConsumerGroupOption) *ConsumerGroup {
	opts := newDefaultConsumerGroupOptions()
	for _, opt := range globalConsumerGroupOptions {
		opt.apply(&opts)
	}
	for _, opt := range options {
		opt.apply(&opts)
	}
	return &ConsumerGroup{
		name:    opts.Name,
		members: opts.Members,
		mut:     sync.RWMutex{},
		pel:     make(PendingEntriesList),
		startAt: opts.StartAt,
	}
}

// GetStartAt returns the start at entry ID for the Consumer group.
func (c *ConsumerGroup) GetStartAt() EntryID {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.startAt
}

// SetStartAt sets the start at entry ID for the Consumer group.
func (c *ConsumerGroup) SetStartAt(id EntryID) {
	c.mut.Lock()
	c.startAt = id
	c.mut.Unlock()
}

// GetName returns the name of the Consumer group.
func (c *ConsumerGroup) GetName() string {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.name
}

// AddMember adds a Consumer group member to the Consumer group. If a Consumer group member with the same name already
// exists, this function overwrites it.
func (c *ConsumerGroup) AddMember(member Consumer) {
	c.mut.Lock()
	c.members[member.name] = member
	c.mut.Unlock()
}

// RemoveMember removes the Consumer group member with the given name. If the member does not exist, this function does
// nothing.
func (c *ConsumerGroup) RemoveMember(member Consumer) {
	c.mut.Lock()
	delete(c.members, member.name)
	c.mut.Unlock()
}

// ListMembers returns a list of all Consumer group members.
func (c *ConsumerGroup) ListMembers() []Consumer {
	c.mut.RLock()
	members := make([]Consumer, 0, len(c.members))
	for _, m := range c.members {
		members = append(members, m)
	}
	c.mut.RUnlock()
	return members
}

// GetMember returns the Consumer group member with the given name. If the member does not exist, this function returns
// false.
func (c *ConsumerGroup) GetMember(name string) (*Consumer, bool) {
	c.mut.RLock()
	m, ok := c.members[name]
	c.mut.RUnlock()
	return &m, ok
}

// GetPendingEntry returns the pending entry with the given ID from the Consumer group's Pending Entries List. If the
// pending entry does not exist, this function returns false.
func (c *ConsumerGroup) GetPendingEntry(id EntryID) (PendingEntry, bool) {
	c.mut.RLock()
	pe, ok := c.pel[id]
	c.mut.RUnlock()
	return pe, ok
}

// GetPendingEntriesForConsumer returns all pending entries for the Consumer group member with the given name.
func (c *ConsumerGroup) GetPendingEntriesForConsumer(consumer string) []PendingEntry {
	c.mut.RLock()
	defer c.mut.RUnlock()
	var out []PendingEntry
	for _, pe := range c.pel {
		if pe.Consumer == consumer {
			out = append(out, pe)
		}
	}
	return out
}

// AddPendingEntry adds a pending entry to the Consumer group's Pending Entries List. The pending entry is associated
// with the given ID and Consumer.
func (c *ConsumerGroup) AddPendingEntry(id EntryID, consumer string) {
	c.mut.Lock()
	c.pel[id] = PendingEntry{
		ID:            id,
		Consumer:      consumer,
		DeliveredAt:   time.Now(),
		DeliveryCount: 1,
	}
	c.mut.Unlock()
}

// incrementDeliveryCountAndTime increments the delivery count and sets the DeliveredAt time for the pending entry with
// the given ID. If the pending entry does not exist, this function does nothing.
func (c *ConsumerGroup) incrementDeliveryCountAndTime(id EntryID) {
	c.mut.Lock()
	pe, ok := c.pel[id]
	if !ok {
		c.mut.Unlock()
		return
	}
	pe.DeliveryCount++
	pe.DeliveredAt = time.Now()
	c.pel[id] = pe
	c.mut.Unlock()
}

// RemovePendingEntry removes the pending entry with the given ID from the Consumer group's Pending Entries List. If the
// pending entry does not exist, this function does nothing.
func (c *ConsumerGroup) RemovePendingEntry(id EntryID) {
	c.mut.Lock()
	delete(c.pel, id)
	c.mut.Unlock()
}

// GetPendingEntries returns all pending entries in the Consumer group's Pending Entries List.
//
// This method returns a copy of the PendingEntriesList. The caller is free to modify the returned list without
// affecting the Consumer group's Pending Entries List.
func (c *ConsumerGroup) GetPendingEntries() PendingEntriesList {
	c.mut.RLock()
	defer c.mut.RUnlock()
	out := make(PendingEntriesList, len(c.pel))
	for id, pe := range c.pel {
		out[id] = pe
	}
	return out
}
