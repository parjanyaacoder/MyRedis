package core

import (
	"strconv"
)

func deduceTypeEncoding(value string) (uint8, uint8) {
	oType := OBJ_TYPE_STRING

	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return oType, OBJ_ENCODING_INT
	}
	if len(value) <= 44 {
		return oType, OBJ_ENCODING_EMBSTR
	}
	return oType, OBJ_ENCODING_RAW
}
