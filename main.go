package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"ChirpyBlog/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecretKey   string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)

	platform := os.Getenv("PLATFORM")

	jwtSecretKey := os.Getenv("SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("SECRET_KEY environment variable is not set")
	}

	const filepathRoot = "."
	const port = ":8080"
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecretKey:   jwtSecretKey,
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	//mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpByID)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetMetrics)

	fmt.Println("Server is running on http://localhost:8080")

	log.Fatal(srv.ListenAndServe())
}
