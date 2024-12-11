// Package historitor was created to provide a transactional log with the following features:
//   - Search by log entry payload
//   - Allow re-writing of log entries (but not deleting)
//   - Allow for multiple readers and writers
//   - Allow for grouping of readers (akin to Kafka Consumer groups)
//   - Backed by persistent storage
//   - Expiration of read group members
//
// The package is heavily inspired by Kafka and Redis Streams.
//
// # Pending Entries List (PEL)
//
// Every Consumer group keeps track of the log entries that have been delivered to its members. This allows us to
// distribute log entries among the members of the Consumer group and ensure that each log entry is processed by only
// one member of the group.
//
// This feature is implemented using a Pending Entries List (PEL) associated with each Consumer group. The PEL is a
// list of log entries that have been delivered to the Consumer group but have not yet been acknowledged by the
// Consumer. The PEL contains information on when the entry was delivered to the Consumer, the number of times the
// entry has been delivered, and the Consumer that received the entry.
//
// # Handling busy consumers
//
// Every entry read from the log must be acknowledged by the Consumer. As entries are read, they are added to the
// Pending Entries List (PEL) for the Consumer group. The PEL includes information on when the entry was delivered to
// the Consumer, the number of times the entry has been delivered, and the Consumer that received the entry.
//
// When a Consumer requests a [Log.Read] operation, the log will check the Pending Entries List (PEL) to see if the
// Consumer has any entries that have not been acknowledged, is older than [WithLogMaxPendingAge], or has been delivered
// more than [WithLogMaxDeliveryCount] times. If the Consumer has any such entries, the log will update the PEL and
// include the entries in the response.
//
// To prevent a Consumer from holding onto an entry indefinitely, a housekeeping function called [Log.Cleanup] is
// implemented. This function, among other things, removes pending entries that have been delivered more than
// [WithLogMaxDeliveryCount] times and are older than [WithLogAttemptRedeliveryAfter].
//
// # Handling dead consumers
//
// A dead Consumer is a Consumer that stops consuming log entries. This can happen for a variety of reasons, such as
// network issues, the Consumer crashing, or the Consumer being shut down. Dead consumers can cause log entries to
// accumulate in the Pending Entries List (PEL) and never be processed. To handle dead consumers, the log implements
// a housekeeping function called [Log.Cleanup]. Among other things, this function removed pending entries that are
// older than [WithLogMaxPendingAge] to allow other consumers to attempt to process the log entry.
package historitor
