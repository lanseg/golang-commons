package common

import (
	"encoding/json"
	"os"
)

// FromJsonFile decodes json file into an object
func FromJsonFile[T any](path string) (*T, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FromJson[T](bytes)
}

// FromJson just tries to convert json string to object
func FromJson[T any](data []byte) (*T, error) {
	result := new(T)
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	return result, nil
}
