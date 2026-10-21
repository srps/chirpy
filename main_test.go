package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func TestReadinessRoute(t *testing.T) {
	response := httptest.NewRecorder()
	apiRouter(&apiConfig{}).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("readiness status = %d, want 200", response.Code)
	}
}

func TestRedirectSlashesDoesNotUseUntrustedHost(t *testing.T) {
	router := chi.NewRouter()
	router.Use(middleware.RedirectSlashes)
	router.Get("/resource", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	request := httptest.NewRequest(http.MethodGet, "http://localhost/resource/", nil)
	request.Host = "attacker.example"
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if location := response.Header().Get("Location"); location != "/resource" {
		t.Fatalf("redirect location = %q, want relative /resource", location)
	}
}
