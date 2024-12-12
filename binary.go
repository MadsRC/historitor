package historitor

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

// An unsigned integer is sent one of two ways. If it is less than 128, it is sent as a byte with that value. Otherwise
// it is sent as a minimal-length big-endian (high byte first) byte stream holding the value, preceded by one byte
// holding the byte count, negated. Thus 0 is transmitted as (00), 7 is transmitted as (07) and 256 is transmitted as
// (FE 01 00).
func encodeUnsignedInt(x uint64) []byte {
	if x <= 0x7F {
		return []byte{uint8(x)}
	}

	b := make([]byte, 9)
	binary.BigEndian.PutUint64(b[1:], x)
	bc := bits.LeadingZeros64(x) >> 3 // 8 - bytelen(x)
	b[bc] = uint8(bc - 8)             // and then we subtract 8 (size of uint64) to get -bytelen(x)

	return b[bc : 8+1]
}

func decodeUnsignedInt(b []byte) (uint64, error) {
	if b[0] <= 0x7f {
		return uint64(b[0]), nil
	}
	n := -int(int8(b[0]))
	if n > 8 {
		return 0, fmt.Errorf("invalid uint data length %d", n)
	}
	if len(b) < n {
		fmt.Errorf("invalid uint data length %d: exceeds input size %d", n, len(b))
	}

	var x uint64
	// Don't need to check error; it's safe to loop regardless.
	// Could check that the high byte is zero but it's not worth it.
	for _, byt := range b[1 : n+1] {
		x = x<<8 | uint64(byt)
	}
	return x, nil
}
