package main

import "net/http"

func main() {
	ServeMux := http.NewServeMux()
	ServeMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	ServeMux.Handle("/assets/logo.png", http.FileServer(http.Dir(".")))
	ServeMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	MyServer := http.Server{}
	MyServer.Handler = ServeMux
	MyServer.Addr = ":8080"
	MyServer.ListenAndServe()
}
