package utils

import (
	"strconv"
	"strings"
)

func ParseHexUint64(s string) (uint64, error) {
	s = strings.TrimPrefix(strings.ToLower(s), "0x")
	return strconv.ParseUint(s, 16, 64)
}
