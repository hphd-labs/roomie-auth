package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
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
	if r.ContentLength > 0 {
		contentTypes := ParseContentTypes(r.Header.Get("Content-Type"))
		if contentTypes[0].Type != "application/json" {
			return false
		}
	}

	// make sure JSON is in the accepts
	for _, acceptType := range ParseContentTypes(r.Header.Get("Accept")) {
		if acceptType.Type == "application/json" {
			return true
		}
	}

	return false
}

type ContentType struct {
	Type    string
	Options map[string]string
}

func ParseContentTypes(data string) []*ContentType {
	types := make([]*ContentType, 0, 1)

	for _, entry := range strings.Split(data, ",") {
		t := &ContentType{
			Options: make(map[string]string),
		}

		components := strings.Split(entry, ";")
		t.Type = components[0]

		for _, opt := range components[1:] {
			values := strings.Split(opt, "=")
			if len(values) != 2 {
				continue
			}
			key := strings.TrimSpace(values[0])
			t.Options[key] = strings.TrimSpace(values[1])
		}

		types = append(types, t)
	}
	return types
}
