// Package orderedmap provides a Map that preserves the order of key value pairs.
package orderedmap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// Map implements a map that keeps the item order.
// It preserves the insertion order of keys for both JSON unmarshalling
// and manual operations.
type Map struct {
	// Data holds the map entries with their sequence indices.
	Data map[string]Entry
	// Keys contains the map keys in insertion order.
	Keys []string
}

// Len returns the number of elements within the map.
func (m *Map) Len() int {
	return len(m.Keys)
}

// Get retrieves the value for the given key.
// Returns the value and true if the key exists, or nil and false otherwise.
func (m *Map) Get(key string) (any, bool) {
	if m.Data == nil {
		return nil, false
	}
	entry, exists := m.Data[key]
	if !exists {
		return nil, false
	}
	return entry.Value, true
}

// Set adds or updates a key-value pair in the map.
// If the key already exists, its value is updated without changing position.
// If the key is new, it is appended to the end of the insertion order.
func (m *Map) Set(key string, value any) {
	if m.Data == nil {
		m.Data = make(map[string]Entry)
	}

	if _, exists := m.Data[key]; !exists {
		m.Keys = append(m.Keys, key)
	}

	m.Data[key] = Entry{
		index: nextSequence(),
		Value: value,
	}
}

// Delete removes a key-value pair from the map.
// Returns true if the key existed and was deleted, false otherwise.
//
// Performance note: This operation is O(n) where n is the number of keys,
// as it requires searching through the Keys slice to find and remove the key.
func (m *Map) Delete(key string) bool {
	if m.Data == nil {
		return false
	}

	if _, exists := m.Data[key]; !exists {
		return false
	}

	delete(m.Data, key)

	for i, k := range m.Keys {
		if k == key {
			m.Keys = append(m.Keys[:i], m.Keys[i+1:]...)
			break
		}
	}

	return true
}

// GetAndDelete retrieves and removes a key in one operation.
// Returns the value and true if the key existed, or nil and false otherwise.
//
// Performance note: This operation is O(n) due to the underlying Delete operation.
func (m *Map) GetAndDelete(key string) (value any, loaded bool) {
	value, loaded = m.Get(key)
	if loaded {
		m.Delete(key)
	}
	return value, loaded
}

// Clear removes all entries from the map, resulting in an empty map.
func (m *Map) Clear() {
	m.Data = make(map[string]Entry)
	m.Keys = nil
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (m *Map) Range(f func(key string, value any) bool) {
	if m.Data == nil || len(m.Keys) == 0 {
		return
	}

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
		return fmt.Errorf("failed to unmarshal orderedmap: %w", err)
	}

	m.rebuildKeys()
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (m *Map) MarshalJSON() ([]byte, error) {
	if m == nil || len(m.Keys) == 0 {
		return []byte("{}"), nil
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(m.Keys)*20))
	buf.WriteString("{")

	lastIdx := len(m.Keys) - 1
	for i, key := range m.Keys {
		keyData, err := json.Marshal(key)
		if err != nil {
			return nil, fmt.Errorf("marshalling key %q: %w", key, err)
		}
		buf.Write(keyData)
		buf.WriteByte(':')

		value := m.Data[key].Value
		b, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("marshalling entry for key %q: %w", key, err)
		}
		buf.Write(b)

		if i < lastIdx {
			buf.WriteByte(',')
		}
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// rebuildKeys builds the keys slice sorted by insertion order.
func (m *Map) rebuildKeys() {
	m.Keys = make([]string, 0, len(m.Data))
	for name := range m.Data {
		m.Keys = append(m.Keys, name)
	}

	sort.SliceStable(m.Keys, func(i, j int) bool {
		return m.Data[m.Keys[i]].index < m.Data[m.Keys[j]].index
	})
}
