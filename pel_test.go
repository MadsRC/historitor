package historitor

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPendingEntriesList_JSON(t *testing.T) {
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
