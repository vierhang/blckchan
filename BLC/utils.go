package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

// JSONToSlice 标准JSON格式转切片
func JSONToSlice(jsonString string) []string {
	var strSlice []string
	if err := json.Unmarshal([]byte(jsonString), &strSlice); err != nil {
		log.Panicf("json to []string failed! %v\n", err)
	}
	return strSlice
}
