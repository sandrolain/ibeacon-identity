package main

import (
	"encoding/binary"
	"hash/crc32"
)

func checkSignature(ib *IBeacon, secret []byte) bool {
	data := append(ib.UUID, secret...)
	crc32Checksum := crc32.ChecksumIEEE(data)

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, crc32Checksum)

	major := binary.BigEndian.Uint16(bs[0:2])
	minor := binary.BigEndian.Uint16(bs[2:4])

	return ib.Major == major && ib.Minor == minor
}
