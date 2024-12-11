package historitor

import (
	"fmt"
	"time"
)

// PendingEntriesList keeps track of log entries that have been delivered to a Consumer group member but not yet
// acknowledged.
type PendingEntriesList map[EntryID]PendingEntry

// String returns a string representation of the PendingEntriesList.
func (pel PendingEntriesList) String() string {
	var out []byte
	for id, entry := range pel {
		out = append(out, id.String()...)
		out = append(out, ':')
		out = append(out, fmt.Sprintf("\n\tConsumer: %s\n\tDelivered at: %s\n\tDelivery count: %d", entry.Consumer, entry.DeliveredAt.UTC(), entry.DeliveryCount)...)
		out = append(out, '\n')
	}

	return string(out)
}

// PendingEntry is an entry in the PendingEntriesList. It keeps track of the Consumer group member that the log entry
// was delivered to, the time it was delivered, and the number of times it has been delivered.
type PendingEntry struct {
	// ID of the log entry
	ID EntryID `json:"ID"`
	// Name of the Consumer group member
	Consumer string `json:"Consumer"`
	// Time the Entry was delivered to the Consumer group member
	DeliveredAt time.Time `json:"delivered_at"`
	// The number of times the Entry has been delivered to the Consumer group member
	DeliveryCount int `json:"delivery_count"`
}
