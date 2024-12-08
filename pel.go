package historitor

import "time"

type pendingEventsList map[string]pendingEvent

type pendingEvent struct {
	// id of the log entry
	id EntryID
	// Name of the consumer group member
	consumer string
	// Time the event was delivered to the consumer group member
	deliveredAt time.Time
	// The number of times the event has been delivered to the consumer group member
	deliveryCount int
}
