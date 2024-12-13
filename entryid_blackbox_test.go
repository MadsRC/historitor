//go:build integration

package historitor_test

import (
	"github.com/MadsRC/historitor"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEntryID_IsZero(t *testing.T) {
	e := historitor.ZeroEntryID
	require.True(t, e.IsZero())
	e = historitor.EntryID{}
	require.True(t, e.IsZero())
	e = historitor.NewEntryID(time.UnixMilli(1733946462442), 10)
	require.False(t, e.IsZero())
}

func TestEntryID_String(t *testing.T) {
	e := historitor.NewEntryID(time.UnixMilli(1733946462442), 10)
	require.Equal(t, "1733946462442-0000000000010", e.String())
}

func TestEntryID_MarshalJSON(t *testing.T) {
	e := historitor.NewEntryID(time.UnixMilli(1733946462442), 10)
	b, err := e.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "\"1733946462442-0000000000010\"", string(b))
}

func TestEntryID_UnmarshalJSON(t *testing.T) {
	var e historitor.EntryID
	expected := historitor.NewEntryID(time.UnixMilli(1733946462442), 10)
	err := e.UnmarshalJSON([]byte("\"1733946462442-0000000000010\""))
	require.NoError(t, err)
	require.Equal(t, expected, e)
}

func TestEntryID_UnmarshalJSON_bad_value(t *testing.T) {
	var e historitor.EntryID
	err := e.UnmarshalJSON([]byte("bad value"))
	require.Error(t, err)
}
