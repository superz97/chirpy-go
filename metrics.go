package main

import (
	"chirpy-go/internal/database"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// @Router /admin/metrics [get]
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(
		"<html>\n"+
			"<body>\n    "+
			"<h1>Welcome, Chirpy Admin</h1>\n    "+
			"<p>Chirpy has been visited %d times!</p>\n  "+
			"</body>\n"+
			"</html>",
		cfg.fileserverHits.Load())))
}

// @Router /admin/reset [post]
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Reset is only allowed in dev environment")
		return
	}
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
