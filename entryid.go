package historitor

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ZeroEntryID        = EntryID{}
	StartFromBeginning = EntryID{
		Time: ZeroEntryID.Time,
		Seq:  128,
	}
	StartFromEnd = EntryID{
		Time: ZeroEntryID.Time,
		Seq:  129,
	}
)

type EntryID struct {
	Time time.Time
	Seq  uint64
}

func (e EntryID) String() string {
	return fmt.Sprintf("%d-%013d", e.Time.UTC().UnixMilli(), e.Seq)
}

func NewEntryID(s string) (EntryID, error) {
	var e EntryID
	ms, err := strconv.Atoi(strings.Split(s, "-")[0])
	if err != nil {
		return ZeroEntryID, fmt.Errorf("failed to parse EntryID: %w", err)
	}
	t := time.UnixMilli(int64(ms))
	seq, err := strconv.Atoi(strings.Split(s, "-")[1])
	if err != nil {
		return ZeroEntryID, fmt.Errorf("failed to parse EntryID: %w", err)
	}
	e.Time = t
	e.Seq = uint64(seq)
	return e, nil
}

func (e EntryID) IsZero() bool {
	return e == ZeroEntryID
}
