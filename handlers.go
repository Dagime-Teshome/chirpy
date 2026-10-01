package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Dagime-Teshome/chirpy/internal/auth"
	"github.com/Dagime-Teshome/chirpy/internal/database"
	"github.com/google/uuid"
)

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}
func (cfg *apiConfig) returnCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", " text/html;charset=utf-8")
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
	w.Header().Add("Content-Type", " text/plain;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Count reset"))
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", " text/plain;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) Handle_createUser(w http.ResponseWriter, r *http.Request) {
	var req_body = user_create{}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	err = json.Unmarshal(data, &req_body)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	pass_hash, err := auth.HashPassword(req_body.Password)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	user_param := database.CreateUserParams{
		Email:          req_body.Email,
		HashedPassword: pass_hash,
	}
	db_user, err := cfg.queries.CreateUser(r.Context(), user_param)
	if err != nil {
		respondWithError(w, 500, "couldn't create user")
		return
	}
	user_json := user_response{
		Email:      db_user.Email,
		Created_at: db_user.CreatedAt,
		Updated_at: db_user.UpdatedAt,
		ID:         db_user.ID.String(),
	}
	err = respondWithJSON(w, 201, user_json)

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
	// parsed_user_id, _ := uuid.Parse(Req_body.User_id)
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	user_id, err := auth.ValidateJWT(token, cfg.secret)
	fmt.Println("id " + user_id.String())
	if err != nil {
		respondWithError(w, 401, "Unauthorized ")
		return
	}
	chirp_param := database.CreateChirpParams{
		UserID: user_id,
		Body:   cleanText,
	}
	db_chipr, err := cfg.queries.CreateChirp(r.Context(), chirp_param)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
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

func (cfg *apiConfig) Handle_Login(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	var user_data = user_create{}
	err = json.Unmarshal(data, &user_data)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	db_user, err := cfg.queries.GerUserByEmail(r.Context(), user_data.Email)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	match, err := auth.CheckPasswordHash(user_data.Password, db_user.HashedPassword)

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	if !match {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	token, err := auth.MakeJWT(db_user.ID, cfg.secret)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	ref_tok_args := database.Create_Refresh_TokenParams{
		Token:  auth.MakeRefreshToken(),
		UserID: db_user.ID,
	}
	db_ref_token, err := cfg.queries.Create_Refresh_Token(r.Context(), ref_tok_args)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	user_json := user_response{
		Email:         db_user.Email,
		Created_at:    db_user.CreatedAt,
		Updated_at:    db_user.UpdatedAt,
		ID:            db_user.ID.String(),
		Token:         token,
		Refresh_Token: db_ref_token.Token,
	}
	err = respondWithJSON(w, 200, user_json)
}

func (cfg *apiConfig) Handle_Token_Revoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	err = cfg.queries.Revoke_Refresh_Token(r.Context(), token)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 204, "Toke Revoked")
}
func (cfg *apiConfig) Handle_Token_Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	db_rf_tok, err := cfg.queries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	if db_rf_tok.RevokedAt.Valid || !db_rf_tok.ExpiresAt.After(time.Now()) {
		respondWithError(w, 401, "refresh token invalid")
		return
	}
	jwt_token, err := auth.MakeJWT(db_rf_tok.UserID, cfg.secret)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	respondWithJSON(w, 201, token_refresh{Token: jwt_token})
}
