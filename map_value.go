package orderedmap

import (
	"encoding/json"
	"fmt"
)

// Entry represents a map value entry with its insertion sequence.
type Entry struct {
	index uint64
	Value any
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *Entry) UnmarshalJSON(b []byte) error {
	e.index = nextSequence()

	if err := json.Unmarshal(b, &e.Value); err != nil {
		return fmt.Errorf("unmarshalling entry: %w", err)
	}

	if _, ok := e.Value.(map[string]any); ok {
		var m Map
		if err := json.Unmarshal(b, &m); err != nil {
			return fmt.Errorf("unmarshalling json entry into map: %w", err)
		}
		e.Value = &m
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (e *Entry) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(e.Value)
	if err != nil {
		return nil, fmt.Errorf("marshalling entry: %w", err)
	}
	return b, nil
}
