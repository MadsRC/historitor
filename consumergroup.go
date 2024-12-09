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
	pel     pendingEntriesList
	startAt EntryID
}

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
		pel:     make(pendingEntriesList),
		startAt: opts.StartAt,
	}
}
func (c *ConsumerGroup) GetStartAt() EntryID {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.startAt
}

func (c *ConsumerGroup) SetStartAt(id EntryID) {
	c.mut.Lock()
	c.startAt = id
	c.mut.Unlock()
}

func (c *ConsumerGroup) GetName() string {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.name
}

func (c *ConsumerGroup) addMember(member Consumer) {
	c.mut.Lock()
	c.members[member.name] = member
	c.mut.Unlock()
}

func (c *ConsumerGroup) removeMember(member Consumer) {
	c.mut.Lock()
	delete(c.members, member.name)
	c.mut.Unlock()
}

func (c *ConsumerGroup) listMembers() []Consumer {
	c.mut.RLock()
	members := make([]Consumer, 0, len(c.members))
	for _, m := range c.members {
		members = append(members, m)
	}
	c.mut.RUnlock()
	return members
}

func (c *ConsumerGroup) getMember(name string) (*Consumer, bool) {
	c.mut.RLock()
	m, ok := c.members[name]
	c.mut.RUnlock()
	return &m, ok
}

func (c *ConsumerGroup) getPendingEntry(id EntryID) (*pendingEntry, bool) {
	c.mut.RLock()
	pe, ok := c.pel[id]
	c.mut.RUnlock()
	return &pe, ok
}

func (c *ConsumerGroup) getPendingEntriesForConsumer(consumer string) []pendingEntry {
	c.mut.RLock()
	defer c.mut.RUnlock()
	var out []pendingEntry
	for _, pe := range c.pel {
		if pe.consumer == consumer {
			out = append(out, pe)
		}
	}
	return out
}

func (c *ConsumerGroup) addPendingEntry(id EntryID, consumer string) {
	c.mut.Lock()
	c.pel[id] = pendingEntry{
		id:            id,
		consumer:      consumer,
		deliveredAt:   time.Now(),
		deliveryCount: 1,
	}
	c.mut.Unlock()
}

// incrementDeliveryCountAndTime increments the delivery count and sets the deliveredAt time for the pending entry with
// the given ID. If the pending entry does not exist, this function does nothing.
func (c *ConsumerGroup) incrementDeliveryCountAndTime(id EntryID) {
	c.mut.Lock()
	pe, ok := c.pel[id]
	if !ok {
		c.mut.Unlock()
		return
	}
	pe.deliveryCount++
	pe.deliveredAt = time.Now()
	c.pel[id] = pe
	c.mut.Unlock()
}

func (c *ConsumerGroup) removePendingEntry(id EntryID) {
	c.mut.Lock()
	delete(c.pel, id)
	c.mut.Unlock()
}
