package imdl

import "errors"

func encodeIndex(index uint32, bytes []byte, byteIndex int) error {
	if byteIndex+2 >= len(bytes) {
		return errors.New("index erorr")
	}
	bytes[byteIndex+0] = byte(index & 0x000000ff)
	bytes[byteIndex+1] = byte((index & 0x0000ff00) >> 8)
	bytes[byteIndex+2] = byte((index & 0x00ff0000) >> 16)
	return nil
}

func decodeIndex(index int, bytes []byte) uint32 {
	byteIndex := index * 3
	return uint32(bytes[byteIndex]) | uint32(bytes[byteIndex+1])<<8 | uint32(bytes[byteIndex+2])<<16
}

func EncodeVertexIndices(indices []uint32) []byte {
	len := len(indices) * 3
	data := make([]byte, len)
	for i := range indices {
		encodeIndex(indices[i], data, i*3)
	}
	return data
}

func DecodeVertexIndices(bytes []byte) []uint32 {
	len := len(bytes) / 3
	inds := make([]uint32, len)
	for i := range inds {
		inds[i] = decodeIndex(i, bytes)
	}
	return inds
}
