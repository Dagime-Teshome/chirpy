package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Dagime-Teshome/chirpy/internal/database"
	"github.com/google/uuid"
)

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}
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

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", " text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) Handle_createUser(w http.ResponseWriter, r *http.Request) {
	var email_body = user_create{}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "couldn't marsahll data")
		return
	}
	json.Unmarshal(data, &email_body)
	db_user, err := cfg.queries.CreateUser(r.Context(), email_body.Email)
	if err != nil {
		respondWithError(w, 500, "couldn't create user")
		return
	}

	err = respondWithJSON(w, 201, db_user)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

}
func (cfg *apiConfig) Handle_DeleteUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		respondWithError(w, 403, "not allowed")
		return
	}
	err := cfg.queries.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, chirp_body{Body: "users delted succesfully"})
}
func (cfg *apiConfig) Handle_Chirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var Req_body = chirp_body{}
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
	parsed_user_id, _ := uuid.Parse(Req_body.User_id)
	chirp_param := database.CreateChirpParams{
		UserID: parsed_user_id,
		Body:   cleanText,
	}
	db_chipr, err := cfg.queries.CreateChirp(r.Context(), chirp_param)
	if err != nil {
		respondWithError(w, 500, err.Error())
	}
	respondWithJSON(w, 200, db_chipr)
}
func (cfg *apiConfig) List_Chirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.queries.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, chirps)

}

func (cfg *apiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	string_id := r.PathValue("chirpID")
	id, err := uuid.Parse(string_id)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	chirp, err := cfg.queries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}
	respondWithJSON(w, 200, chirp)
}
