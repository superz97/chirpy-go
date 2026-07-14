package main

import (
	"log"
	"net/http"
)

const port = "8080"

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/", fileServer)

	server := &http.Server{}
	server.Handler = mux
	server.Addr = ":" + port

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
