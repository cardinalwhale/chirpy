package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	apiCfg := apiConfig{}
	ServeMux := http.NewServeMux()
	ServeMux.Handle("/app", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	ServeMux.Handle("/assets/logo.png", http.FileServer(http.Dir(".")))
	ServeMux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	ServeMux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		fmt.Fprintf(w, "Hits: %d", apiCfg.fileServerHits.Load())
	})
	ServeMux.HandleFunc("POST /reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileServerHits.Swap(0)
	})
	MyServer := http.Server{}
	MyServer.Handler = ServeMux
	MyServer.Addr = ":8080"
	MyServer.ListenAndServe()
}
