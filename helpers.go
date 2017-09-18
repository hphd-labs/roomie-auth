package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// DecodeJSON wraps the creation of a decoder and a quick decode
func DecodeJSON(data io.Reader, obj interface{}) error {
	return json.NewDecoder(data).Decode(obj)
}

// RenderJSON wraps the creation of an encoder and a quick encode,
// as well as setting necessary headers in the reply
func RenderJSON(w http.ResponseWriter, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(obj)
}

// SupportsJSON checks for bi-directional JSON encoding support
func AcceptsJSON(r *http.Request) bool {
	// if there's content, make sure it's JSON
	if r.ContentLength > 0 && r.Header.Get("Content-Type") != "application/json" {
		return false
	}

	// make sure JSON is in the accepts
	for _, acceptEncoding := range r.Header["Accept"] {
		if acceptEncoding == "application/json" {
			return true
		}
	}

	return false
}
