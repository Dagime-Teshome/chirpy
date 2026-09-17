package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, 400, map[string]string{"error": msg})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {

	dat, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(dat)
	return nil
}
