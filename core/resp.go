package core

import (
	"errors"
	"fmt"
)

func readLength(data []byte) (int, int) {
	pos, length := 0, 0
	for pos = range data {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			return length, pos + 2
		}
		length = length*10 + int(b-'0')
	}
	return 0, 0
}

func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for ; data[pos] != '\r'; pos++ {

	}
	return string(data[1:pos]), pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readInt64(data []byte) (int64, int, error) {
	pos := 1
	var val int64 = 0

	for ; data[pos] != '\r'; pos++ {
		val = val*10 + int64(data[pos]-'0')
	}

	return val, pos + 2, nil
}

func readBulkString(data []byte) (string, int, error) {
	pos := 1
	len, delta := readLength(data[pos:])
	pos += delta

	return string(data[pos : pos+len]), pos + len + 2, nil
}

func readArray(data []byte) ([]interface{}, int, error) {
	pos := 1
	elementsCount, delta := readLength(data[pos:])
	pos += delta

	var elements []interface{} = make([]interface{}, elementsCount)

	for i := range elements {
		element, delta, error := DecodeOne(data[pos:])
		if error != nil {
			return nil, 0, error
		}

		elements[i] = element
		pos += delta
	}
	return elements, pos, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("No data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}
	return nil, 0, nil
}

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("No data")
	}

	value, _, err := DecodeOne(data)
	return value, err
}

func DecodeArrayStrings(data []byte) ([]string, error) {
	value, err := Decode(data)

	if err != nil {
		return nil, err
	}

	ts := value.([]interface{})
	tokens := make([]string, len(ts))

	for i := range tokens {
		tokens[i] = ts[i].(string)
	}

	return tokens, nil
}

func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		{
			if isSimpleString {
				return []byte(fmt.Sprintf("+%s\r\n", v))
			}
			return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
		}
	case int, int8, int16, int32, int64:
		return []byte(fmt.Sprintf(":%d\r\n", v))
	}
	return []byte{}
}
