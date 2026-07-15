package main

import (
	"log"
	"net/http"
)

const port = "8080"

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app/", fileServer))

	server := &http.Server{}
	server.Handler = mux
	server.Addr = ":" + port

	mux.HandleFunc("/healthz", handleReadiness)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
