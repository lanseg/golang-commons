package almostio

import (
	"encoding/json"
)

type Marshaller[T any] struct {
	Marshal   func(*T) ([]byte, error)
	Unmarshal func([]byte) (*T, error)
}

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
