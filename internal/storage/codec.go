package storage

import "encoding/json"

func encode(value any) ([]byte, error) { return json.Marshal(value) }

func decode(data []byte, value any) error { return json.Unmarshal(data, value) }

func cloneBytes(value []byte) []byte {
	result := make([]byte, len(value))
	copy(result, value)
	return result
}
