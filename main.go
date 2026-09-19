package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/Dagime-Teshome/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        database.Queries
}

func (cfg *apiConfig) returnCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", " text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<html>
						<body>
							<h1>Welcome, Chirpy Admin</h1>
							<p>Chirpy has been visited %d times!</p>
						</body>
					</html>`,
		cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetCount(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", " text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Count reset"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", " text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
func handleChirpValidate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var Req_body = Req_body{}
	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "couldn't marsahll data")
	}
	json.Unmarshal(dat, &Req_body)
	if len(Req_body.Body) >= 140 {
		respondWithError(w, 400, "chirp too long")
		return
	}
	replacer := strings.NewReplacer(
		"kerfuffle", "****",
		"Kerfuffle", "****",
		"KERFUFFLE", "****",

		"sharbert", "****",
		"Sharbert", "****",
		"SHARBERT", "****",

		"fornax", "****",
		"Fornax", "****",
		"FORNAX", "****",
	)
	cleanText := replacer.Replace(Req_body.Body)
	respondWithJSON(w, 200, clean_body{Cleaned_Body: cleanText})
}
func main() {
	// load enviroment variablse to os env file
	godotenv.Load()
	// get connection string from os env file
	conn_string := os.Getenv("DB_URL")
	// open a tcp connection to database server
	db, err := sql.Open("postgres", conn_string)
	if err != nil {
		fmt.Println("database connection failed:", err)
		return

	}
	// use sqlc generated code as ORM to add ,remove ,update data from database
	database_queries := database.New(db)
	// add the sqlc orm thingy to the config so handler and middleware have access to it.
	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		queries:        *database_queries,
	}
	// server config code
	serv_mux := http.NewServeMux()
	serv_mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	// api end points
	serv_mux.HandleFunc("GET /api/healthz", handleHealthz)
	serv_mux.HandleFunc("POST /api/validate_chirp", handleChirpValidate)
	serv_mux.Handle("GET /admin/metrics", http.HandlerFunc(apiConfig.returnCount))
	serv_mux.Handle("POST /admin/reset", http.HandlerFunc(apiConfig.resetCount))
	// server code
	var srv http.Server
	srv.Addr = ":8080"
	srv.Handler = serv_mux
	srv.ListenAndServe()
}
