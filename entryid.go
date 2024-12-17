package historitor

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ZeroEntryID        = EntryID{}
	StartFromBeginning = EntryID{
		time: ZeroEntryID.time,
		seq:  128,
	}
	StartFromEnd = EntryID{
		time: ZeroEntryID.time,
		seq:  129,
	}
)

// EntryID is a unique identifier for an entry in a log.
type EntryID struct {
	time time.Time
	seq  uint64
}

// NewEntryID creates a new EntryID with the given time and sequence number.
// The time is truncated to milliseconds and timezone set to UTC.
func NewEntryID(t time.Time, seq uint64) EntryID {
	return EntryID{
		time: t.Truncate(time.Millisecond).UTC(),
		seq:  seq,
	}
}

func (e EntryID) IsZero() bool {
	return e == ZeroEntryID
}

func (e EntryID) String() string {
	return fmt.Sprintf("%d-%013d", e.time.UTC().UnixMilli(), e.seq)
}

func (e EntryID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", e.String())), nil
}

func (e *EntryID) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	id, err := ParseEntryID(s)
	if err != nil {
		return fmt.Errorf("failed to unmarshal EntryID: %w", err)
	}
	*e = id
	return nil
}

func (e EntryID) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(externalEntryID{
		Time: e.time,
		Seq:  e.seq,
	})
	return buf.Bytes(), nil
}

func (e *EntryID) UnmarshalBinary(data []byte) error {
	var ee externalEntryID
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&ee)
	if err != nil {
		return err
	}
	e.time = ee.Time
	e.seq = ee.Seq
	return nil
}

// ParseEntryID parses a string representation of an EntryID.
// The string must be in the format "time-seq" where time is the number of milliseconds since the Unix epoch and seq is
// the sequence number. The time is truncated to milliseconds and timezone set to UTC.
func ParseEntryID(s string) (EntryID, error) {
	var e EntryID
	ms, err := strconv.Atoi(strings.Split(s, "-")[0])
	if err != nil {
		return ZeroEntryID, fmt.Errorf("failed to parse EntryID: %w", err)
	}
	t := time.UnixMilli(int64(ms)).UTC()
	seq, err := strconv.Atoi(strings.Split(s, "-")[1])
	if err != nil {
		return ZeroEntryID, fmt.Errorf("failed to parse EntryID: %w", err)
	}
	e.time = t
	e.seq = uint64(seq)
	return e, nil
}

// externalEntryID is used to represent a version of en [EntryID] that can easily be encoded and decoded by the gob
// package
type externalEntryID struct {
	Time time.Time
	Seq  uint64
}
