# orderedmap

[![Build status](https://github.com/cornelk/orderedmap/actions/workflows/go.yaml/badge.svg?branch=main)](https://github.com/cornelk/orderedmap/actions)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/cornelk/orderedmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/cornelk/orderedmap)](https://goreportcard.com/report/github.com/cornelk/orderedmap)
[![codecov](https://codecov.io/gh/cornelk/orderedmap/branch/main/graph/badge.svg?token=NS5UY28V3A)](https://codecov.io/gh/cornelk/orderedmap)

## Overview

A Golang Map that keeps track of the insert order of items.

The current version can be used as a replacement for `map[string]interface{}` and allows JSON unmarshalling into it.
No manual adding of items is currently possible.

## Usage

Unmarshall and marshall JSON:

```
var m Map
input := `{"423":"abc","231":"dbh","152":"xyz"}`

m.UnmarshalJSON([]byte(input))

output, _ := m.MarshalJSON()
# unlike the standard Golang map, the output will be exact the same and not in random key order
```
