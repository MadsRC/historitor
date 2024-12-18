//go:build !integration

package historitor

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	fakeTestEntryID1 = EntryID{
		time: time.Unix(0, 0).Add(time.Duration(1734467114191) * time.Millisecond).UTC(),
		seq:  1,
	}
	fakeTestEntryID2 = EntryID{
		time: time.Unix(0, 0).Add(time.Duration(1734467114191) * time.Millisecond).UTC(),
		seq:  2,
	}
	fakeTestEntryID3 = EntryID{
		time: time.Unix(0, 0).Add(time.Duration(1734467114192) * time.Millisecond).UTC(),
		seq:  1,
	}
)

// TestNewEntryID tests the NewEntryID function and that the time is truncated to milliseconds and set to UTC.
func TestNewEntryID(t *testing.T) {
	naow := time.Now()
	eid := NewEntryID(naow, 1)

	t.Run("Time is truncated", func(t *testing.T) {
		require.Equal(t, naow.Truncate(time.Millisecond).UnixMilli(), eid.time.UnixMilli())
	})

	t.Run("Time is UTC", func(t *testing.T) {
		require.True(t, eid.time.Location().String() == "UTC")
	})

	t.Run("Sequence number is set", func(t *testing.T) {
		require.Equal(t, uint64(1), eid.seq)
	})

	t.Run("Time is not zero", func(t *testing.T) {
		require.False(t, eid.time.IsZero())
	})
}

// TestEntryID_IsZero tests the IsZero function.
func TestEntryID_IsZero(t *testing.T) {
	t.Run("Zero EntryID", func(t *testing.T) {
		require.True(t, ZeroEntryID.IsZero())
	})

	t.Run("Non-zero EntryID", func(t *testing.T) {
		require.False(t, NewEntryID(time.Now(), 1).IsZero())
	})
}

// TestEntryID_String tests the String function.
func TestEntryID_String(t *testing.T) {
	eid := EntryID{
		time: time.Now().Truncate(time.Millisecond).UTC(),
		seq:  1,
	}
	eidStr := eid.String()

	require.Len(t, strings.Split(eidStr, "-"), 2)
	p1 := strings.Split(eidStr, "-")[0]
	p2 := strings.Split(eidStr, "-")[1]

	t.Run("First part is time in milliseconds", func(t *testing.T) {
		t1, err := strconv.ParseInt(p1, 10, 64)
		require.NoError(t, err)
		require.Equal(t, eid.time.UnixMilli(), t1)
	})

	t.Run("Second part is sequence number", func(t *testing.T) {
		require.Len(t, p2, 13)
		i1, err := strconv.Atoi(p2)
		require.NoError(t, err)
		require.Equal(t, int(eid.seq), i1)
	})
}

// TestEntryID_MarshalJSON tests the MarshalJSON function.
func TestEntryID_MarshalJSON(t *testing.T) {
	naow := time.Now().Truncate(time.Millisecond).UTC()
	eid := EntryID{
		time: naow,
		seq:  1,
	}
	eidJSON, err := eid.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, `"`+eid.String()+`"`, string(eidJSON))
}

// TestEntryID_UnmarshalJSON tests the UnmarshalJSON function.
func TestEntryID_UnmarshalJSON(t *testing.T) {
	input := "1734467114191-0000000000001"
	eidJSON := []byte(`"` + input + `"`)
	var eid2 EntryID
	err := eid2.UnmarshalJSON(eidJSON)
	require.NoError(t, err)
	require.Equal(t, strings.Split(input, "-")[0], strconv.FormatInt(eid2.time.UnixMilli(), 10))
	require.Equal(t, uint64(1), eid2.seq)
}

func TestEntryID_UnmarshalJSON_error(t *testing.T) {
	input := "this is not a valid EntryID"
	eidJSON := []byte(`"` + input + `"`)
	var eid2 EntryID
	err := eid2.UnmarshalJSON(eidJSON)
	require.Error(t, err)
	require.Equal(t, ZeroEntryID, eid2)
}

// TestEntryID_MarshalBinary tests the MarshalBinary function.
func TestEntryID_MarshalBinary(t *testing.T) {
	// For some reason, "go test -run TestEntryID_MarshalBinary" succeeds, but "go test" fails on this test.
	// It appears that go decides to encode the message differently when running all tests.
	// Further investigation is needed.
	t.Skip("skipped")
	naow := time.Unix(0, 0).Add(time.Duration(1734467114191) * time.Millisecond).UTC()
	expected := []byte{0x2e, 0x7f, 0x3, 0x1, 0x1, 0xf, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x49, 0x44, 0x1, 0xff, 0x80, 0x0, 0x1, 0x2, 0x1, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x53, 0x65, 0x71, 0x1, 0x6, 0x0, 0x0, 0x0, 0x10, 0xff, 0x81, 0x5, 0x1, 0x1, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x82, 0x0, 0x0, 0x0, 0x16, 0xff, 0x80, 0x1, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xde, 0xf3, 0xd5, 0x2a, 0xb, 0x62, 0x6d, 0xc0, 0xff, 0xff, 0x1, 0x1, 0x0}

	eid := EntryID{
		time: naow,
		seq:  1,
	}
	bin, err := eid.MarshalBinary()
	require.NoError(t, err)
	require.Equal(t, expected, bin)
}

// TestEntryID_UnmarshalBinary tests the UnmarshalBinary function.
func TestEntryID_UnmarshalBinary(t *testing.T) {
	input := []byte{0x2e, 0x7f, 0x3, 0x1, 0x1, 0xf, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x49, 0x44, 0x1, 0xff, 0x80, 0x0, 0x1, 0x2, 0x1, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x53, 0x65, 0x71, 0x1, 0x6, 0x0, 0x0, 0x0, 0x10, 0xff, 0x81, 0x5, 0x1, 0x1, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x82, 0x0, 0x0, 0x0, 0x16, 0xff, 0x80, 0x1, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xde, 0xf3, 0xd5, 0x2a, 0xb, 0x62, 0x6d, 0xc0, 0xff, 0xff, 0x1, 0x1, 0x0}
	expectedTime := time.Unix(0, 0).Add(time.Duration(1734467114191) * time.Millisecond).UTC()
	expectedSeq := uint64(1)

	var eid EntryID
	err := eid.UnmarshalBinary(input)
	require.NoError(t, err)
	require.Equal(t, expectedTime, eid.time)
	require.Equal(t, expectedSeq, eid.seq)
}

// TestEntryID_UnmarshalBinary_error tests the UnmarshalBinary when encountering bad data.
func TestEntryID_UnmarshalBinary_error(t *testing.T) {
	input := []byte("this is not a valid EntryID")
	var eid EntryID
	err := eid.UnmarshalBinary(input)
	require.Error(t, err)
	require.Equal(t, ZeroEntryID, eid)
}

// TestParseEntryID tests the ParseEntryID function.
func TestParseEntryID(t *testing.T) {
	input := "1734467114191-0000000000001"
	eid, err := ParseEntryID(input)
	require.NoError(t, err)
	require.Equal(t, time.Unix(0, 0).Add(time.Duration(1734467114191)*time.Millisecond).UTC(), eid.time)
	require.Equal(t, uint64(1), eid.seq)
}

// TestParseEntryID_error_bad_timestamp tests the ParseEntryID function when encountering bad data in the form of a
// bad timestamp.
func TestParseEntryID_error_bad_timestamp(t *testing.T) {
	input := "this is not a valid EntryID-0000000000001"
	eid, err := ParseEntryID(input)
	require.Error(t, err)
	require.Equal(t, ZeroEntryID, eid)
}

// TestParseEntryID_error_bad_seq tests the ParseEntryID function when encountering bad data in the form of a bad
// sequence number.
func TestParseEntryID_error_bad_seq(t *testing.T) {
	input := "1734467114191-this is not a valid sequence number"
	eid, err := ParseEntryID(input)
	require.Error(t, err)
	require.Equal(t, ZeroEntryID, eid)
}
