package queue

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

// Serialize serializes a Queue into bytes using cbor
func Serialize[T any](queue *Queue[T]) ([]byte, error) {
	queue.Lock()
	defer queue.Unlock()

	elems := make([]T, 0, queue.queue.Len())
	next := queue.queue.Front()
	for next != nil {
		elems = append(elems, next.Value.(T))
		next = next.Next()
	}

	data, err := cbor.Marshal(elems)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Deserialize deserializes cbor data into a Queue
func Deserialize[T any](data []byte) (*Queue[T], error) {
	var elems []T
	if err := decMode.Unmarshal(data, &elems); err != nil {
		return nil, err
	}

	queue := New[T]()
	queue.Push(elems...)
	return queue, nil
}
