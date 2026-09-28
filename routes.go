package main

import "net/http"

func routes(config *apiConfig) *http.ServeMux {
	serv_mux := http.NewServeMux()
	serv_mux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serv_mux.Handle("GET /admin/metrics", http.HandlerFunc(config.returnCount))
	serv_mux.HandleFunc("POST /admin/reset", config.Handle_DeleteUsers)
	// api end points
	serv_mux.HandleFunc("GET /api/chirps/{chirpID}", config.GetChirp)
	serv_mux.HandleFunc("GET /api/chirps", config.List_Chirps)
	serv_mux.HandleFunc("POST /api/chirps", config.Handle_Chirp)
	serv_mux.HandleFunc("POST /api/users", config.Handle_createUser)
	serv_mux.HandleFunc("GET /api/healthz", handleHealthz)

	return serv_mux
}
