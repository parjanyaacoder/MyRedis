func DecodeOne(data []byte) 

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("No data")
	}

	value, _ , err := DecodeOne(data);
	return value, err;
}