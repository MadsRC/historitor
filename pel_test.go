//go:build !integration

package historitor

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPendingEntriesList_MarshalJSON(t *testing.T) {
	pel := make(PendingEntriesList)
	ti := time.Now().Truncate(time.Millisecond).UTC()
	id := EntryID{time: ti, seq: 0}
	pel[id] = PendingEntry{
		ID:            id,
		Consumer:      "consumer1",
		DeliveredAt:   ti,
		DeliveryCount: 1,
	}

	_, err := json.Marshal(pel)
	require.NoError(t, err)
	_, err = json.Marshal(pel[id])
	require.NoError(t, err)
}

func TestPendingEntriesList_String(t *testing.T) {
	expected := "1734467114191-0000000000000:\n\tConsumer: consumer1\n\tDelivered at: 2024-12-17 20:25:14.191 +0000 UTC\n\tDelivery count: 1\n"
	ti := time.Unix(0, 0).Add(time.Duration(1734467114191) * time.Millisecond).UTC()
	pel := make(PendingEntriesList)
	id := EntryID{time: ti, seq: 0}
	pel[id] = PendingEntry{
		ID:            id,
		Consumer:      "consumer1",
		DeliveredAt:   ti,
		DeliveryCount: 1,
	}

	require.Equal(t, expected, pel.String())
}

func TestPendingEntry_String(t *testing.T) {
	expected := "1734467114191-0000000000000:\n\tConsumer: consumer1\n\tDelivered at: 2024-12-17 20:25:14.191 +0000 UTC\n\tDelivery count: 1"
	ti := time.Unix(0, 0).Add(time.Duration(1734467114191) * time.Millisecond).UTC()
	id := EntryID{time: ti, seq: 0}
	pe := PendingEntry{
		ID:            id,
		Consumer:      "consumer1",
		DeliveredAt:   ti,
		DeliveryCount: 1,
	}

	require.Equal(t, expected, pe.String())
}
