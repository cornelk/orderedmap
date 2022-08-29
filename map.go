// Package orderedmap provides a Map that preserves the order of key value pairs.
package orderedmap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// Map implements a map that keeps the item order.
type Map struct {
	// Data holds the map entries in unsorted order.
	Data map[string]Entry
	// Keys contains the map keys in sorted order.
	Keys []string
}

// Len returns the number of elements within the map.
func (m *Map) Len() int {
	return len(m.Keys)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (m *Map) Range(f func(key string, value interface{}) bool) {
	for _, key := range m.Keys {
		entry := m.Data[key]
		if !f(key, entry.Value) {
			return
		}
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Map) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &m.Data); err != nil {
		return err
	}

	m.rebuildKeys()
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (m *Map) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")

	for i, key := range m.Keys {
		buf.WriteString(fmt.Sprintf("%q:", key))

		value := m.Data[key].Value
		b, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("marshalling entry: %w", err)
		}
		buf.Write(b)

		if i < len(m.Keys)-1 {
			buf.WriteString(",")
		}
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}

// rebuildKeys build the sorted keys slice.
func (m *Map) rebuildKeys() {
	m.Keys = []string{}
	for name := range m.Data {
		m.Keys = append(m.Keys, name)
	}

	sort.SliceStable(m.Keys, func(i, j int) bool {
		return m.Data[m.Keys[i]].index < m.Data[m.Keys[j]].index
	})
}
