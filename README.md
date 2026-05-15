# orderedmap

[![Build status](https://github.com/cornelk/orderedmap/actions/workflows/go.yaml/badge.svg?branch=main)](https://github.com/cornelk/orderedmap/actions)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/cornelk/orderedmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/cornelk/orderedmap)](https://goreportcard.com/report/github.com/cornelk/orderedmap)
[![codecov](https://codecov.io/gh/cornelk/orderedmap/branch/main/graph/badge.svg?token=NS5UY28V3A)](https://codecov.io/gh/cornelk/orderedmap)

## Overview

A Go map that preserves insertion order for manual operations and JSON marshalling.
Nested JSON objects are unmarshalled as ordered maps too.

## Installation

```bash
go get github.com/cornelk/orderedmap
```

Requires Go 1.24 or later.

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/cornelk/orderedmap"
)

func main() {
	var m orderedmap.Map

	m.Set("name", "Alice")
	m.Set("age", 30)

	if value, ok := m.Get("name"); ok {
		fmt.Println(value)
	}

	m.Range(func(key string, value any) bool {
		fmt.Printf("%s: %v\n", key, value)
		return true
	})

	data, _ := json.Marshal(&m)
	fmt.Println(string(data))
}
```

JSON input order is preserved when unmarshalling and marshalling:

```go
var m orderedmap.Map
input := `{"423":"abc","231":"dbh","152":"xyz"}`

json.Unmarshal([]byte(input), &m)
output, _ := json.Marshal(&m)
// Output: {"423":"abc","231":"dbh","152":"xyz"}
```

## API

- `Set(key string, value any)` adds or updates a value without moving existing keys.
- `Get(key string) (any, bool)` retrieves a value.
- `Delete(key string) bool` removes a value.
- `GetAndDelete(key string) (any, bool)` retrieves and removes a value.
- `Clear()` removes all values.
- `Len() int` returns the number of entries.
- `Range(func(key string, value any) bool)` iterates in insertion order.

## License

Apache-2.0 License - See LICENSE file for details.
