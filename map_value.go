package orderedmap

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Entry represents a map value entry.
type Entry struct {
	index uint64
	Value interface{}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *Entry) UnmarshalJSON(b []byte) error {
	s.index = nextSequence()

	if err := json.Unmarshal(b, &s.Value); err != nil {
		return fmt.Errorf("unmarshalling entry: %w", err)
	}

	if s.Value != nil {
		valueKind := reflect.TypeOf(s.Value).Kind()
		if valueKind == reflect.Map { // force values of type map to be of ordered map as well
			var m Map
			if err := json.Unmarshal(b, &m); err != nil {
				return fmt.Errorf("unmarshalling json entry into map: %w", err)
			}
			s.Value = m
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (s *Entry) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}
