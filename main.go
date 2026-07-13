package main

import "net/http"

func main() {
	ServeMux := http.NewServeMux()
	MyServer := http.Server{}
	MyServer.Handler = ServeMux
	MyServer.Addr = ":8080"
	MyServer.ListenAndServe()
}
