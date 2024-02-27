package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

func main() {

	addr := "localhost:8080"
	r := chi.NewRouter()
	apiCfg := &apiConfig{fileServerHits: 0}
	publicFSHandler := NewFileServerHandler("public/")
	appPrefix := "/app"
	wrappedFSHandler := apiCfg.middlewareMetricsInc(http.StripPrefix(appPrefix, publicFSHandler))
	r.Handle(appPrefix, wrappedFSHandler)
	r.Handle(appPrefix+"/*", wrappedFSHandler)
	r.Mount("/api", apiRouter(apiCfg))
	r.Mount("/admin", adminRouter(apiCfg))
	wrappedMux := NewCorsMiddleware(r)
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

func apiRouter(apiCfg *apiConfig) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", readinessHandler)
	r.HandleFunc("/reset", apiCfg.ResetServerHits)
	return r
}

func adminRouter(apiCfg *apiConfig) http.Handler {
	r := chi.NewRouter()
	r.Get("/metrics", apiCfg.RenderServerHits)
	return r
}

func NewFileServerHandler(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}
