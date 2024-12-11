package historitor

import (
	"encoding/json"
	"fmt"
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

	b, err := json.Marshal(pel)
	require.NoError(t, err)
	fmt.Println(string(b))
	b, err = json.Marshal(pel[id])
	require.NoError(t, err)
	fmt.Println(string(b))
}
