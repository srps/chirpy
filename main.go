package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type apiConfig struct {
	fileServerHits int
}

func main() {

	addr := "localhost:8080"
	mux := http.NewServeMux()
	appPrefix := "/app/"
	apiCfg := &apiConfig{fileServerHits: 0}
	fileServerHandler := http.HandlerFunc(handleFileServing)
	mux.Handle(appPrefix, http.StripPrefix(appPrefix, apiCfg.middlewareMetricsInc(fileServerHandler)))
	mux.HandleFunc("/healthz", readinessHandler)
	mux.HandleFunc("/metrics", apiCfg.fileServerHitsHandler)
	mux.HandleFunc("/reset", apiCfg.ResetServerHits)
	wrappedMux := NewCorsMiddleware(mux)
	srv := &http.Server{
		Addr:         addr,
		Handler:      wrappedMux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("server is listening at %s", addr)
	log.Fatal(srv.ListenAndServe())
}

func handleFileServing(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("public/")).ServeHTTP(w, r)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits += 1
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) fileServerHitsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits)))
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