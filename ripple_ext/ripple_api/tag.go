package ripple_api

import (
	"strconv"
)

func BuildTag(str string) (*uint32, error) {
	if str == "" {
		return nil, nil
	}

	tag := new(uint32)

	value, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return nil, err
	}

	*tag = uint32(value)

	return tag, nil
}

func ParseTag(tag *uint32) string {
	if tag == nil {
		return ""
	}

	return strconv.FormatUint(uint64(*tag), 10)
}