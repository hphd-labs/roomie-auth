package main

import (
	"encoding/json"
	content "github.com/andrewburian/content-type"
	"io"
	"net/http"
)

var contentTypeJson *content.ContentType

func init() {
	contentTypeJson, _ = content.ParseSingle("application/json")
}

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

	ct, accepts, err := content.ParseRequest(r)
	if err != nil {
		return false
	}

	// if there's content, make sure it's JSON
	if r.ContentLength > 0 {
		if ct.MediaType != contentTypeJson.MediaType {
			return false
		}
	}

	// make sure JSON is in the accepts
	return accepts.SupportsType(contentTypeJson)
}
