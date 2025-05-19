package types

import "encoding/json"

type Error struct {
	Code    string                     `json:"code"`
	Message string                     `json:"message"`
	Extra   map[string]json.RawMessage `json:"extra,omitempty"`
}
