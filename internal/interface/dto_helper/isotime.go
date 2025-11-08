package dto_helper

import (
	"silo/pkg/database/seedwork"
	"silo/pkg/types/isotime"
)

// ISOTime .
func ISOTime(in isotime.ISOTime) isotime.ISOTime {
	return in
}

// ISOTimeToUnix .
func ISOTimeToUnix(in isotime.ISOTime) int64 {
	if in.IsZero() {
		return 0
	}
	return in.ToTime().Unix()
}

func WithTimestamp(in seedwork.WithTimestamp) seedwork.WithTimestamp {
	return in
}
