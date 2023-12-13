package common

import (
	"encoding/json"
)

// FromJson just tries to convert json string to object
func FromJson[T any](data []byte) (*T, error) {
	result := new(T)
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	return result, nil
}
