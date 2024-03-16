package tjson

import (
	"encoding/json"
)

func Unmarshal[T any](data []byte) (v T, err error) {
	return v, json.Unmarshal(data, &v)
}