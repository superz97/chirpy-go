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
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	apiCfg := &apiConfig{
		db:       dbQueries,
		platform: platform,
	}

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fileServer)))

	server := &http.Server{}
	server.Handler = mux
	server.Addr = ":" + port

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
