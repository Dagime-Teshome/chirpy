package main

import (
	"sync/atomic"
	"time"

	"github.com/Dagime-Teshome/chirpy/internal/database"
)

type chirp_body struct {
	Body    string `json:"body"`
	User_id string `json:"user_id"`
}

type Res_body struct {
	Valid bool `json:"valid"`
}

type clean_body struct {
	Cleaned_Body string `json:"cleaned_body"`
}

type err_resp struct {
	Error string `json:"error"`
}

type user_create struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type user_response struct {
	ID            string    `json:"id"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Token         string    `json:"token,omitempty"`
	Refresh_Token string    `json:"refresh_token,omitempty"`
	Is_Chirpy_Red bool      `json:"is_chirpy_red "`
}

type apiConfig struct {
	Platform       string
	fileserverHits atomic.Int32
	queries        database.Queries
	secret         string
	api_key        string
}
type token_refresh struct {
	Token string `json:"token"`
}
type apiHandler struct{}

type hook_data struct {
	User_id string `json:"user_id"`
}
type hook_body struct {
	Event string    `json:"event"`
	Data  hook_data `json:"data"`
}
