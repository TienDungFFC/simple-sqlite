package main

import (
	"bytes"
)

func readVarInt(buf []byte) (result uint64, count int) {
	result = 0
	count = 0
	shiff := 0
	for _, b := range buf {
		result |= uint64(b&0x7F) << shiff
		count++
		if b&0x80 == 0 {
			return
		}
		shiff += 7
	}
	return
}

func readVarint(reader *bytes.Reader) (result uint64, byteRead int) {
	var shift uint
	result = 0
	byteRead = 0
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return
		}
		byteRead++
		result |= uint64(b&0x7F) << shift

		if b&0x80 == 0 {
			break
		}
		shift += 7
		if shift >= 64 {
			return
		}
	}

	return
}
