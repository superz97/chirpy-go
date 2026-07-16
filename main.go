package main

import (
	"log"
	"net/http"
)

const port = "8080"

func main() {
	mux := http.NewServeMux()

	apiCfg := &apiConfig{}

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fileServer)))

	server := &http.Server{}
	server.Handler = mux
	server.Addr = ":" + port

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /api/reset", apiCfg.handlerReset)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
