package router

import (
	"github.com/gorilla/mux"
	"news-api/internal/http/handlers"
)

func NewRouter(authHandler *handlers.AuthHandler) *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	// Auth-user routing
	api.HandleFunc("/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")

	r.Use(mux.CORSMethodMiddleware(r))

	return r
}
