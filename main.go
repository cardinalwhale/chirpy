package main

import (
	"encoding/json"
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
	ServeMux.Handle("/app/assets/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	ServeMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	ServeMux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("<html> <body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html>", apiCfg.fileServerHits.Load())))
	})
	ServeMux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileServerHits.Swap(0)
	})
	ServeMux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		type parameter struct {
			Body string `json:"body"`
		}

		decoder := json.NewDecoder(r.Body)
		param := parameter{}
		err := decoder.Decode(&param)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if len(param.Body) <= 140 {
			type validReturnParam struct {
				ValidReturn bool `json:"valid"`
			}
			validReturn := validReturnParam{ValidReturn: true}
			dat, err := json.Marshal(validReturn)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(dat)
		} else {
			type errorReturnParam struct {
				ErrorReturn string `json:"error"`
			}
			errorReturn := errorReturnParam{ErrorReturn: "Chirp is too long"}
			dat, err := json.Marshal(errorReturn)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write(dat)
		}
	})
	MyServer := http.Server{}
	MyServer.Handler = ServeMux
	MyServer.Addr = ":8080"
	MyServer.ListenAndServe()
}
