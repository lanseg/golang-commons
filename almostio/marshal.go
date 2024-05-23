package almostio

import (
	"encoding/json"
)

// Marshaller contains reader and writer for a given format.
// Used when you need to write same data into different file types (e.g. json, binary, text...)
type Marshaller[T any] struct {
	Marshal   func(*T) ([]byte, error)
	Unmarshal func([]byte) (*T, error)
}

// NewJsonMarshal creates a reader/writer for Json file type
func NewJsonMarshal[T any]() *Marshaller[T] {
	return &Marshaller[T]{
		Marshal: func(obj *T) ([]byte, error) {
			return json.Marshal(obj)
		},
		Unmarshal: func(data []byte) (*T, error) {
			t := new(T)
			err := json.Unmarshal(data, t)
			return t, err
		},
	}
}
