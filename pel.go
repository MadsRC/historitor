package historitor

import "time"

type pendingEntriesList map[EntryID]pendingEntry

type pendingEntry struct {
	// id of the log entry
	id EntryID
	// Name of the consumer group member
	consumer string
	// Time the Entry was delivered to the consumer group member
	deliveredAt time.Time
	// The number of times the Entry has been delivered to the consumer group member
	deliveryCount int
}
