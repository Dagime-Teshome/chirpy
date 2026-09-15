package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) returnCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", " text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", cfg.fileserverHits.Load())
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

func handleHealthz (w http.ResponseWriter , r *http.Request) {
		w.Header().Add("Content-Type", " text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
}
}
func main() {
	apiConfig := apiConfig{}
	serv_mux := http.NewServeMux()
	serv_mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	// api end points
	serv_mux.HandleFunc("GET api/healthz", handleHealthz)
	serv_mux.Handle("GET api/metrics", http.HandlerFunc(apiConfig.returnCount))
	serv_mux.Handle("POST api/reset", http.HandlerFunc(apiConfig.resetCount))
	// server code
	var srv http.Server
	srv.Addr = ":8080"
	srv.Handler = serv_mux
	srv.ListenAndServe()
}
