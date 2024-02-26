package main

import (
	"net/http"
)

type CorsWrapper struct {
	next http.Handler
}

func NewCorsMiddleware(next http.Handler) *CorsWrapper {
	return &CorsWrapper{next: next}
}

func (h *CorsWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	h.next.ServeHTTP(w, r)
}
