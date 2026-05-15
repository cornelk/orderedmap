package orderedmap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestJSON(t *testing.T) {
	input := `{"423":"abc","231":"dbh","152":"xyz"}`

	var m Map
	if err := m.UnmarshalJSON([]byte(input)); err != nil {
		t.Fatal(err)
	}

	output, err := m.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	got := string(output)
	if input != got {
		t.Errorf("expected: %v, got: %v", input, got)
	}
}

func TestJSON_EmptyMap(t *testing.T) {
	var m Map
	output, err := m.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	expected := "{}"
	got := string(output)
	if expected != got {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestJSON_EscapesInvalidUTF8Keys(t *testing.T) {
	var m Map
	m.Set("bad"+string([]byte{0xff})+"key", "value")

	output, err := json.Marshal(&m)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(output) {
		t.Fatalf("expected valid JSON, got: %s", output)
	}

	var decoded map[string]string
	if err := json.Unmarshal(output, &decoded); err != nil {
		t.Fatal(err)
	}
	if got := decoded["bad\ufffdkey"]; got != "value" {
		t.Errorf("expected replacement-character key to contain value, got: %q", got)
	}
}

func TestJSON_NestedMap(t *testing.T) {
	input := `{"outer":{"inner":"value"},"key":"val"}`

	var m Map
	if err := m.UnmarshalJSON([]byte(input)); err != nil {
		t.Fatal(err)
	}

	outer, ok := m.Data["outer"].Value.(*Map)
	if !ok {
		t.Fatalf("expected nested object to unmarshal as *Map, got: %T", m.Data["outer"].Value)
	}
	if got, ok := outer.Get("inner"); !ok || got != "value" {
		t.Fatalf("expected nested value inner=value, got: %v, %v", got, ok)
	}

	output, err := m.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	got := string(output)
	if input != got {
		t.Errorf("expected: %v, got: %v", input, got)
	}
}

func TestLen(t *testing.T) {
	input := `{"423":"abc","231":"dbh","152":"xyz"}`

	var m Map
	if err := m.UnmarshalJSON([]byte(input)); err != nil {
		t.Fatal(err)
	}

	if got, expected := m.Len(), 3; expected != got {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestLen_EmptyMap(t *testing.T) {
	var m Map
	if got := m.Len(); got != 0 {
		t.Errorf("expected: 0, got: %v", got)
	}
}

func TestGet_Success(t *testing.T) {
	var m Map
	m.Set("key1", "value1")
	m.Set("key2", 42)

	val, ok := m.Get("key1")
	if !ok {
		t.Error("expected key1 to exist")
	}
	if val != "value1" {
		t.Errorf("expected: value1, got: %v", val)
	}

	val, ok = m.Get("key2")
	if !ok {
		t.Error("expected key2 to exist")
	}
	if val != 42 {
		t.Errorf("expected: 42, got: %v", val)
	}
}

func TestGet_NotFound(t *testing.T) {
	var m Map
	m.Set("key1", "value1")

	val, ok := m.Get("nonexistent")
	if ok {
		t.Error("expected key to not exist")
	}
	if val != nil {
		t.Errorf("expected: nil, got: %v", val)
	}
}

func TestGet_EmptyMap(t *testing.T) {
	var m Map
	val, ok := m.Get("key")
	if ok {
		t.Error("expected key to not exist in empty map")
	}
	if val != nil {
		t.Errorf("expected: nil, got: %v", val)
	}
}

func TestSet_NewKeys(t *testing.T) {
	var m Map
	m.Set("first", 1)
	m.Set("second", 2)
	m.Set("third", 3)

	if m.Len() != 3 {
		t.Errorf("expected length: 3, got: %v", m.Len())
	}

	expectedKeys := []string{"first", "second", "third"}
	if !reflect.DeepEqual(m.Keys, expectedKeys) {
		t.Errorf("expected keys: %v, got: %v", expectedKeys, m.Keys)
	}
}

func TestSet_UpdateExisting(t *testing.T) {
	var m Map
	m.Set("key1", "original")
	m.Set("key2", "value2")
	m.Set("key1", "updated")

	if m.Len() != 2 {
		t.Errorf("expected length: 2, got: %v", m.Len())
	}

	val, ok := m.Get("key1")
	if !ok {
		t.Error("expected key1 to exist")
	}
	if val != "updated" {
		t.Errorf("expected: updated, got: %v", val)
	}

	expectedKeys := []string{"key1", "key2"}
	if !reflect.DeepEqual(m.Keys, expectedKeys) {
		t.Errorf("expected keys: %v, got: %v", expectedKeys, m.Keys)
	}
}

func TestDelete_Success(t *testing.T) {
	var m Map
	m.Set("key1", "value1")
	m.Set("key2", "value2")
	m.Set("key3", "value3")

	ok := m.Delete("key2")
	if !ok {
		t.Error("expected delete to succeed")
	}

	if m.Len() != 2 {
		t.Errorf("expected length: 2, got: %v", m.Len())
	}

	_, exists := m.Get("key2")
	if exists {
		t.Error("expected key2 to be deleted")
	}

	expectedKeys := []string{"key1", "key3"}
	if !reflect.DeepEqual(m.Keys, expectedKeys) {
		t.Errorf("expected keys: %v, got: %v", expectedKeys, m.Keys)
	}
}

func TestDelete_NotFound(t *testing.T) {
	var m Map
	m.Set("key1", "value1")

	ok := m.Delete("nonexistent")
	if ok {
		t.Error("expected delete to fail for nonexistent key")
	}

	if m.Len() != 1 {
		t.Errorf("expected length: 1, got: %v", m.Len())
	}
}

func TestDelete_EmptyMap(t *testing.T) {
	var m Map
	ok := m.Delete("key")
	if ok {
		t.Error("expected delete to fail on empty map")
	}
}

func TestRange(t *testing.T) {
	input := `{"423":"abc","231":"dbh","152":"xyz"}`

	var m Map
	if err := m.UnmarshalJSON([]byte(input)); err != nil {
		t.Fatal(err)
	}

	var items []map[string]any
	m.Range(func(key string, value any) bool {
		entry := map[string]any{key: value}
		items = append(items, entry)
		return true
	})

	expected := []map[string]any{
		{"423": "abc"},
		{"231": "dbh"},
		{"152": "xyz"},
	}

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("expected: %v, got: %v", expected, items)
	}
}

func TestRange_EarlyTermination(t *testing.T) {
	var m Map
	m.Set("key1", 1)
	m.Set("key2", 2)
	m.Set("key3", 3)

	var keys []string
	m.Range(func(key string, _ any) bool {
		keys = append(keys, key)
		return len(keys) < 2 // Stop after 2 iterations.
	})

	if len(keys) != 2 {
		t.Errorf("expected 2 iterations, got: %v", len(keys))
	}

	expectedKeys := []string{"key1", "key2"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("expected keys: %v, got: %v", expectedKeys, keys)
	}
}

func TestRange_EmptyMap(t *testing.T) {
	var m Map
	count := 0
	m.Range(func(_ string, _ any) bool {
		count++
		return true
	})

	if count != 0 {
		t.Errorf("expected 0 iterations, got: %v", count)
	}
}

func TestManualAndJSON_Integration(t *testing.T) {
	var m Map
	m.Set("manual1", "value1")
	m.Set("manual2", "value2")

	jsonData, err := m.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var m2 Map
	if err := m2.UnmarshalJSON(jsonData); err != nil {
		t.Fatal(err)
	}

	val, ok := m2.Get("manual1")
	if !ok || val != "value1" {
		t.Errorf("expected manual1=value1, got: %v, %v", val, ok)
	}

	val, ok = m2.Get("manual2")
	if !ok || val != "value2" {
		t.Errorf("expected manual2=value2, got: %v, %v", val, ok)
	}
}

func TestStandardJSONMarshaller(t *testing.T) {
	var m Map
	m.Set("key1", "value1")
	m.Set("key2", 123)

	jsonData, err := json.Marshal(&m)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"key1":"value1","key2":123}`
	got := string(jsonData)
	if expected != got {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestStandardJSONUnmarshaller(t *testing.T) {
	input := []byte(`{"first":"a","second":"b","third":"c"}`)

	var m Map
	if err := json.Unmarshal(input, &m); err != nil {
		t.Fatal(err)
	}

	if m.Len() != 3 {
		t.Errorf("expected length: 3, got: %v", m.Len())
	}

	expectedKeys := []string{"first", "second", "third"}
	if !reflect.DeepEqual(m.Keys, expectedKeys) {
		t.Errorf("expected keys: %v, got: %v", expectedKeys, m.Keys)
	}
}

func TestSet_NilValue(t *testing.T) {
	var m Map
	m.Set("key", nil)

	val, ok := m.Get("key")
	if !ok {
		t.Error("expected key to exist")
	}
	if val != nil {
		t.Errorf("expected nil value, got: %v", val)
	}
}

func TestSet_VariousTypes(t *testing.T) {
	var m Map
	m.Set("string", "value")
	m.Set("int", 42)
	m.Set("float", 3.14)
	m.Set("bool", true)
	m.Set("slice", []int{1, 2, 3})
	m.Set("map", map[string]int{"a": 1})

	if m.Len() != 6 {
		t.Errorf("expected length: 6, got: %v", m.Len())
	}

	tests := []struct {
		key      string
		expected any
	}{
		{"string", "value"},
		{"int", 42},
		{"float", 3.14},
		{"bool", true},
	}

	for _, tt := range tests {
		val, ok := m.Get(tt.key)
		if !ok {
			t.Errorf("expected key %s to exist", tt.key)
		}
		if val != tt.expected {
			t.Errorf("key %s: expected %v, got: %v", tt.key, tt.expected, val)
		}
	}
}

func TestJSON_ComplexNested(t *testing.T) {
	input := `{"level1":{"level2":{"level3":"deep"}}}`

	var m Map
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatal(err)
	}

	output, err := json.Marshal(&m)
	if err != nil {
		t.Fatal(err)
	}

	if input != string(output) {
		t.Errorf("expected: %v, got: %v", input, string(output))
	}
}

func TestJSON_Array(t *testing.T) {
	input := `{"numbers":[1,2,3],"strings":["a","b","c"]}`

	var m Map
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatal(err)
	}

	output, err := json.Marshal(&m)
	if err != nil {
		t.Fatal(err)
	}

	if input != string(output) {
		t.Errorf("expected: %v, got: %v", input, string(output))
	}
}

func TestJSON_MixedTypes(t *testing.T) {
	input := `{"string":"text","number":123,"bool":true,"null":null}`

	var m Map
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatal(err)
	}

	output, err := json.Marshal(&m)
	if err != nil {
		t.Fatal(err)
	}

	if input != string(output) {
		t.Errorf("expected: %v, got: %v", input, string(output))
	}
}

func TestConcurrentRead(_ *testing.T) {
	var m Map
	for i := range 100 {
		m.Set(fmt.Sprintf("key%d", i), i)
	}

	done := make(chan struct{}, 10)
	for range 10 {
		go func() {
			for j := range 100 {
				m.Get(fmt.Sprintf("key%d", j))
			}
			done <- struct{}{}
		}()
	}

	for range 10 {
		<-done
	}
}

func TestLargeMap(t *testing.T) {
	var m Map
	n := 10000

	for i := range n {
		m.Set(fmt.Sprintf("key%d", i), i)
	}

	if m.Len() != n {
		t.Errorf("expected length: %d, got: %d", n, m.Len())
	}

	count := 0
	m.Range(func(key string, _ any) bool {
		expected := fmt.Sprintf("key%d", count)
		if key != expected {
			t.Errorf("expected key: %s, got: %s", expected, key)
		}
		count++
		return true
	})

	if count != n {
		t.Errorf("expected %d iterations, got: %d", n, count)
	}
}

func TestGetAndDelete_ExistingKey(t *testing.T) {
	var m Map
	m.Set("key1", "value1")
	m.Set("key2", "value2")
	m.Set("key3", "value3")

	value, loaded := m.GetAndDelete("key2")
	if !loaded {
		t.Error("expected key2 to exist")
	}
	if value != "value2" {
		t.Errorf("expected value2, got: %v", value)
	}

	if m.Len() != 2 {
		t.Errorf("expected length: 2, got: %d", m.Len())
	}

	_, exists := m.Get("key2")
	if exists {
		t.Error("expected key2 to be deleted")
	}

	expectedKeys := []string{"key1", "key3"}
	if !reflect.DeepEqual(m.Keys, expectedKeys) {
		t.Errorf("expected keys: %v, got: %v", expectedKeys, m.Keys)
	}
}

func TestGetAndDelete_NonExistentKey(t *testing.T) {
	var m Map
	m.Set("key1", "value1")

	value, loaded := m.GetAndDelete("nonexistent")
	if loaded {
		t.Error("expected key to not exist")
	}
	if value != nil {
		t.Errorf("expected nil, got: %v", value)
	}

	if m.Len() != 1 {
		t.Errorf("expected length: 1, got: %d", m.Len())
	}
}

func TestGetAndDelete_EmptyMap(t *testing.T) {
	var m Map

	value, loaded := m.GetAndDelete("key")
	if loaded {
		t.Error("expected key to not exist in empty map")
	}
	if value != nil {
		t.Errorf("expected nil, got: %v", value)
	}
}

func TestClear_WithData(t *testing.T) {
	var m Map
	m.Set("key1", "value1")
	m.Set("key2", "value2")
	m.Set("key3", "value3")

	if m.Len() != 3 {
		t.Errorf("expected length: 3, got: %d", m.Len())
	}

	m.Clear()

	if m.Len() != 0 {
		t.Errorf("expected length: 0, got: %d", m.Len())
	}

	if len(m.Data) != 0 {
		t.Errorf("expected empty Data map, got: %d entries", len(m.Data))
	}

	if m.Keys != nil {
		t.Errorf("expected nil Keys slice, got: %v", m.Keys)
	}

	_, exists := m.Get("key1")
	if exists {
		t.Error("expected all keys to be deleted")
	}
}

func TestClear_EmptyMap(t *testing.T) {
	var m Map

	m.Clear()

	if m.Len() != 0 {
		t.Errorf("expected length: 0, got: %d", m.Len())
	}
}

func TestClear_ThenReuse(t *testing.T) {
	var m Map
	m.Set("old1", "value1")
	m.Set("old2", "value2")

	m.Clear()

	m.Set("new1", "value1")
	m.Set("new2", "value2")

	if m.Len() != 2 {
		t.Errorf("expected length: 2, got: %d", m.Len())
	}

	val, exists := m.Get("new1")
	if !exists || val != "value1" {
		t.Errorf("expected new1=value1, got: %v, %v", val, exists)
	}

	_, exists = m.Get("old1")
	if exists {
		t.Error("expected old keys to not exist after clear and reuse")
	}

	expectedKeys := []string{"new1", "new2"}
	if !reflect.DeepEqual(m.Keys, expectedKeys) {
		t.Errorf("expected keys: %v, got: %v", expectedKeys, m.Keys)
	}
}
