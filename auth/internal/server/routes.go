package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	restLogger "auth/pkg/rest-logger"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()
	r.Use(restLogger.LogMiddleware(slog.Default()))
	r.HandleFunc("/", s.BlankHandler)

	return r
}

func (s *Server) BlankHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "La la la la"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
