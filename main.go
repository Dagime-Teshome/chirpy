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
	// load enviroment variablse to os env file
	godotenv.Load()
	// get connection string from os env file
	conn_string := os.Getenv("DB_URL")
	platform_env := os.Getenv("PLATFORM")
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
		Platform:       platform_env,
		fileserverHits: atomic.Int32{},
		queries:        *database_queries,
	}
	serv_mux := routes(&apiConfig)
	// server code
	var srv http.Server
	srv.Addr = ":8080"
	srv.Handler = serv_mux
	srv.ListenAndServe()
}
