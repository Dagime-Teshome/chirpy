package main

import (
	"time"
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
	Email string `json:"email"`
}
type user_response struct {
	ID         string    `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}
