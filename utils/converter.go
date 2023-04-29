package utils

import "encoding/binary"

func Uint32ToByteArray(val uint32, size int) []byte {
	buf := make([]byte, 4)

	binary.BigEndian.PutUint32(buf[:], val)

	return buf[4-size:]
}

func ByteArrayToUint32(val []byte) uint32 {
	var res uint32
	buf := make([]byte, 4)
	c := 3
	for i := len(val); i > 0; i-- {
		buf[c] = val[i-1]
		c -= 1
	}

	res = binary.BigEndian.Uint32(buf[:])

	return res
}
