package store

import "encoding/json"

func unmarshal(data []byte, target any) error { return json.Unmarshal(data, target) }

func Encode(value any) ([]byte, error) { return json.Marshal(value) }

func Decode(data []byte, target any) error { return json.Unmarshal(data, target) }
