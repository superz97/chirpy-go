package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/", fileServer)

	server := &http.Server{}
	server.Handler = mux
	server.Addr = ":8080"

	log.Fatal(server.ListenAndServe())
}
