package core

import "errors"

func getType(typeEncoding uint8) uint8 {
	return  (typeEncoding >> 4) << 4
}

func getEncoding(typeEncoding uint8) uint8 {
	return  typeEncoding & 0b00001111
}

func assertType(typeEncoding uint8, t uint8) error {
	if getType(typeEncoding) != t {
		return errors.New("the operation is not permitted on this type")
	}
	return nil
}

func assertEncoding(typeEncoding uint8, encoding uint8) error {
	if getEncoding(typeEncoding) != encoding {
		return errors.New("the operation is not permitted on this encoding")
	}
	return nil
}