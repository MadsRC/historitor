package historitor

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

// TestNewConsumerGroup_default_members_different_addresses tests that consumer groups created with NewConsumerGroup
// doesn't share the same members slice address.
//
// We do this using reflection to compare the map as pointers, instead of comparing the map contents or the map itself.
func TestNewConsumerGroup_default_members_different_addresses(t *testing.T) {
	cg1 := NewConsumerGroup(WithConsumerGroupName("cg1"))
	cg2 := NewConsumerGroup(WithConsumerGroupName("cg2"))

	require.False(t, reflect.ValueOf(cg1.members).Pointer() == reflect.ValueOf(cg2.members).Pointer(), "members slices should have different addresses")
}
