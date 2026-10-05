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
	serv_mux.HandleFunc("GET /api/healthz", handleHealthz)
	serv_mux.HandleFunc("POST /api/users", config.Handle_createUser)
	serv_mux.HandleFunc("POST /api/users/login", config.Handle_Login)
	serv_mux.HandleFunc("POST /api/revoke", config.Handle_Token_Revoke)
	serv_mux.HandleFunc("POST /api/refresh", config.Handle_Token_Refresh)
	serv_mux.HandleFunc("PUT /api/users", config.Handle_Update_User)
	serv_mux.HandleFunc("DELETE /api/chirps/{chirpID}", config.Handle_Delete_Chirp)
	serv_mux.HandleFunc("POST /api/polka/webhooks", config.Handle_Hook_Call)

	return serv_mux
}
