package set

import "github.com/fxamacker/cbor/v2"

var decMode cbor.DecMode

func init() {
	mode, err := cbor.DecOptions{
		MaxArrayElements: 300000,
		UTF8:             cbor.UTF8DecodeInvalid,
	}.DecMode()
	if err != nil {
		panic(err)
	}
	decMode = mode
}

// Serialize serializes a HashSet into bytes using cbor
func Serialize[T comparable](set *HashSet[T]) ([]byte, error) {
	set.Lock()
	defer set.Unlock()
	data, err := cbor.Marshal(set.underlying)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Deserialize deserializes cbor data into a HashSet
func Deserialize[T comparable](data []byte) (*HashSet[T], error) {
	underlying := make(map[T]struct{})
	if err := decMode.Unmarshal(data, &underlying); err != nil {
		return nil, err
	}
	return &HashSet[T]{
		underlying: underlying,
	}, nil
}
