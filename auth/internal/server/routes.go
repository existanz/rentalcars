package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"auth/internal/models"
	restLogger "auth/pkg/rest-logger"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()
	r.Use(restLogger.LogMiddleware(slog.Default()))
	r.HandleFunc("/", s.BlankHandler)
	r.HandleFunc("/signup", s.SignUpHandler).Methods("POST")

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

func (s *Server) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var auth models.Auth
	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		log.Fatalf("error unmarshalling JSON. Err: %v", err)
	}
	validate := validator.New()
	err = validate.Struct(auth)
	if err != nil {
		log.Fatalf("error validating JSON. Err: %v", err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("error hashing password. Err: %v", err)
	}
	user := models.User{
		Email:        auth.Email,
		PasswordHash: string(hashedPassword),
	}
	slog.Info("Sign Up Handler", "user", user)
	if err := s.db.RegisterUser(user); err != nil {
		log.Fatalf("error registering user. Err: %v", err)
	}
}
