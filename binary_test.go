//go:build !integration

package historitor

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnsignedIntEncode(t *testing.T) {
	require.Equal(t, []byte{0x00}, encodeUnsignedInt(0))
	require.Equal(t, []byte{0x07}, encodeUnsignedInt(7))
	require.Equal(t, []byte{0xFE, 0x01, 0x00}, encodeUnsignedInt(256))
}

func TestUnsignedIntDecode(t *testing.T) {
	x, err := decodeUnsignedInt([]byte{0x00})
	require.NoError(t, err)
	require.Equal(t, uint64(0), x)

	x, err = decodeUnsignedInt([]byte{0x07})
	require.NoError(t, err)
	require.Equal(t, uint64(7), x)

	x, err = decodeUnsignedInt([]byte{0xFE, 0x01, 0x00})
	require.NoError(t, err)
	require.Equal(t, uint64(256), x)
}
