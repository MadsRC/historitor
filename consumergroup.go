package historitor

import (
	"sync"
	"time"
)

type consumerGroup struct {
	name    string
	members map[string]consumerGroupMember
	mut     sync.RWMutex
	pel     pendingEntriesList
	startAt EntryID
}

func newConsumerGroup(name string, startAt EntryID) consumerGroup {
	return consumerGroup{
		name:    name,
		members: make(map[string]consumerGroupMember),
		mut:     sync.RWMutex{},
		pel:     make(pendingEntriesList),
		startAt: startAt,
	}
}

func (c *consumerGroup) GetStartAt() EntryID {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.startAt
}

func (c *consumerGroup) SetStartAt(id EntryID) {
	c.mut.Lock()
	c.startAt = id
	c.mut.Unlock()
}

func (c *consumerGroup) addMember(member consumerGroupMember) {
	c.mut.Lock()
	c.members[member.name] = member
	c.mut.Unlock()
}

func (c *consumerGroup) removeMember(member consumerGroupMember) {
	c.mut.Lock()
	delete(c.members, member.name)
	c.mut.Unlock()
}

func (c *consumerGroup) listMembers() []consumerGroupMember {
	c.mut.RLock()
	members := make([]consumerGroupMember, 0, len(c.members))
	for _, m := range c.members {
		members = append(members, m)
	}
	c.mut.RUnlock()
	return members
}

func (c *consumerGroup) getMember(name string) (*consumerGroupMember, bool) {
	c.mut.RLock()
	m, ok := c.members[name]
	c.mut.RUnlock()
	return &m, ok
}

func (c *consumerGroup) getPendingEntry(id EntryID) (*pendingEntry, bool) {
	c.mut.RLock()
	pe, ok := c.pel[id]
	c.mut.RUnlock()
	return &pe, ok
}

func (c *consumerGroup) getPendingEntriesForConsumer(consumer string) []pendingEntry {
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

func (c *consumerGroup) addPendingEntry(id EntryID, consumer string) {
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
func (c *consumerGroup) incrementDeliveryCountAndTime(id EntryID) {
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

func (c *consumerGroup) removePendingEntry(id EntryID) {
	c.mut.Lock()
	delete(c.pel, id)
	c.mut.Unlock()
}
