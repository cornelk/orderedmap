package orderedmap

import (
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

func TestRange(t *testing.T) {
	input := `{"423":"abc","231":"dbh","152":"xyz"}`

	var m Map
	if err := m.UnmarshalJSON([]byte(input)); err != nil {
		t.Fatal(err)
	}

	var items []map[string]interface{}
	m.Range(func(key string, value interface{}) bool {
		entry := map[string]interface{}{key: value}
		items = append(items, entry)
		return true
	})

	expected := []map[string]interface{}{
		{"423": "abc"},
		{"231": "dbh"},
		{"152": "xyz"},
	}

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("expected: %v, got: %v", expected, items)
	}
}
