package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
)

// IntToHex 实现int64 -> []byte
func IntToHex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if err != nil {
		log.Panicf("int transact to []byte faild %v\n", err)
	}
	return buffer.Bytes()
}
