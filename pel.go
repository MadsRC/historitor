package historitor

import (
	"encoding/json"
	"fmt"
	"time"
)

// PendingEntriesList keeps track of log entries that have been delivered to a Consumer group member but not yet
// acknowledged.
type PendingEntriesList map[EntryID]PendingEntry

// String returns a string representation of the PendingEntriesList.
func (pel PendingEntriesList) String() string {
	var out []byte
	for _, entry := range pel {
		out = append(out, entry.String()...)
		out = append(out, '\n')
	}

	return string(out)
}

func (pel PendingEntriesList) MarshalJSON() ([]byte, error) {
	type PendingEntry struct {
		Consumer      string    `json:"consumer"`
		DeliveredAt   time.Time `json:"delivered_at"`
		DeliveryCount int       `json:"delivery_count"`
	}
	out := make(map[string]PendingEntry, len(pel))
	for id, entry := range pel {
		out[id.String()] = PendingEntry{
			Consumer:      entry.Consumer,
			DeliveredAt:   entry.DeliveredAt,
			DeliveryCount: entry.DeliveryCount,
		}
	}

	return json.Marshal(out)
}

// PendingEntry is an entry in the PendingEntriesList. It keeps track of the Consumer group member that the log entry
// was delivered to, the time it was delivered, and the number of times it has been delivered.
type PendingEntry struct {
	// ID of the log entry
	ID EntryID `json:"id,omitempty"`
	// Name of the Consumer group member
	Consumer string `json:"consumer"`
	// time the Entry was delivered to the Consumer group member
	DeliveredAt time.Time `json:"delivered_at"`
	// The number of times the Entry has been delivered to the Consumer group member
	DeliveryCount int `json:"delivery_count"`
}

// String returns a string representation of the PendingEntry.
func (pe PendingEntry) String() string {
	return fmt.Sprintf("%s:\n\tConsumer: %s\n\tDelivered at: %s\n\tDelivery count: %d", pe.ID.String(), pe.Consumer, pe.DeliveredAt.UTC(), pe.DeliveryCount)
}
