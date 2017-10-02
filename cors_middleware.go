package main

import (
	"net/http"
	"strings"
)

type CorsHandler struct {
	AllowedMethods []string
	AllowedOrigins []string
	AllowedHeaders []string
}

// Respond to OPTIONS requests
func (h *CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", strings.Join(h.AllowedMethods, ", "))
	w.Header().Add("Access-Control-Allow-Headers", strings.Join(h.AllowedHeaders, ", "))
	w.WriteHeader(http.StatusOK)
}

// Allow origins on all requests that may return a body that needs read
func (h *CorsHandler) ServeHTTPMiddleware(w http.ResponseWriter, r *http.Request,
	n func(http.ResponseWriter, *http.Request)) {

	// Allow all requests access from approved origins
	r.Header.Set("Access-Control-Allow-Origins", strings.Join(h.AllowedOrigins, ", "))
	n(w, r)
}
