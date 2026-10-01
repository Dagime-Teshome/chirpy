package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Dagime-Teshome/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	conn_string := os.Getenv("DB_URL")
	platform_env := os.Getenv("PLATFORM")
	sign_secret := os.Getenv("SECRET")
	db, err := sql.Open("postgres", conn_string)
	if err != nil {
		fmt.Println("database connection failed:", err)
		return

	}
	database_queries := database.New(db)
	apiConfig := apiConfig{
		Platform:       platform_env,
		fileserverHits: atomic.Int32{},
		queries:        *database_queries,
		secret:         sign_secret,
	}
	serv_mux := routes(&apiConfig)
	var srv http.Server
	srv.Addr = ":8080"
	srv.Handler = serv_mux
	srv.ListenAndServe()
}
