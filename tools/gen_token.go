package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"os"

	"github.com/google/uuid"
)

func main() {

	id := uuid.New()

	uuidBytes := id[0:16]

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET environment variable not set")
	}

	data := append(uuidBytes, []byte(secret)...)
	crc32Checksum := crc32.ChecksumIEEE(data)

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, crc32Checksum)

	major := binary.BigEndian.Uint16(bs[0:2])
	minor := binary.BigEndian.Uint16(bs[2:4])

	ummid := append(uuidBytes, bs...)

	fmt.Printf("UMMID:\n%x\n\nUUID:\n%s\n\nMajor:\n%d\n\nMinor:\n%d\n", ummid, id.String(), major, minor)

}
