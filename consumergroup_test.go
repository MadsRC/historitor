package historitor

import (
	"bytes"
	"encoding/gob"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

// TestNewConsumerGroup_default_members_different_addresses tests that Consumer groups created with NewConsumerGroup
// doesn't share the same members slice address.
//
// We do this using reflection to compare the map as pointers, instead of comparing the map contents or the map itself.
func TestNewConsumerGroup_default_members_different_addresses(t *testing.T) {
	cg1 := NewConsumerGroup(WithConsumerGroupName("cg1"))
	cg2 := NewConsumerGroup(WithConsumerGroupName("cg2"))

	require.False(t, reflect.ValueOf(cg1.members).Pointer() == reflect.ValueOf(cg2.members).Pointer(), "members slices should have different addresses")
}

type externalConsumerGroup struct {
	Name    string
	Members map[string]Consumer
	PEL     PendingEntriesList
	StartAt EntryID
}

func (cg *ConsumerGroup) MarshalBinary() ([]byte, error) {
	ecg := externalConsumerGroup{
		Name:    cg.name,
		Members: cg.members,
		PEL:     cg.pel,
		StartAt: cg.startAt,
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(ecg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (cg *ConsumerGroup) UnmarshalBinary(data []byte) error {
	var ecg externalConsumerGroup
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&ecg)
	if err != nil {
		return err
	}
	cg.name = ecg.Name
	cg.members = ecg.Members
	cg.pel = ecg.PEL
	cg.startAt = ecg.StartAt
	return nil
}
