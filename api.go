package main

import (
	"html/template"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

const ContentTypeHtml = "text/html; charset=utf-8"
const templateFile = "templates/metrics.html"

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("file server hit: %s", r.URL.Path)
		cfg.fileServerHits += 1
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) RenderServerHits(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", ContentTypeHtml)
	w.WriteHeader(http.StatusOK)
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Printf("error parsing template: %s", err)
		return
	}
	if err := t.Execute(w, struct{ Count int }{cfg.fileServerHits}); err != nil {
		log.Printf("error executing template: %s", err)
		return
	}
	if err != nil {
		log.Printf("error writing response: %s", err)
		return
	}
	return
}

func (cfg *apiConfig) ResetServerHits(w http.ResponseWriter, _ *http.Request) {
	cfg.fileServerHits = 0
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hits reset to 0"))
	if err != nil {
		log.Printf("error writing response: %s", err)
		return
	}
	return
}
