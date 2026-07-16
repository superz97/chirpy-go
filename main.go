package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"chirpy-go/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const port = "8080"

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()

	apiCfg := &apiConfig{db: dbQueries}

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fileServer)))

	server := &http.Server{}
	server.Handler = mux
	server.Addr = ":" + port

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
