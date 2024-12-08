// Package historitor was created to provide a transactional log with the following features:
// - Search by log entry payload
// - Allow re-writing of log entries (but not deleting)
// - Allow for multiple readers and writers
// - Allow for grouping of readers (akin to Kafka consumer groups)
// - Backed by persistent storage
// - Expiration of read group members
//
// The package is heavily inspired by Kafka and Redis Streams.
package historitor
