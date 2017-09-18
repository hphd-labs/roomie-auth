package main

import (
	"encoding/json"
	"io"
)

// Decode JSON wraps the creation of a decoder and a quick decode
func DecodeJSON(data io.Reader, obj interface{}) error {
	return json.NewDecoder(data).Decode(obj)
}
