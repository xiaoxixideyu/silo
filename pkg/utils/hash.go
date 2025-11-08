package utils

import "hash/fnv"

// HashFnv32 .
func HashFnv32(b []byte) uint32 {
	h := fnv.New32a()
	h.Write(b)
	return h.Sum32()
}

// StrHashToInt32 .
func StrHashToInt32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// StrHashToInt64 .
func StrHashToInt64(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// HashToPartition .
func HashToPartition(s string, partitionCount int) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() % uint32(partitionCount))
}
