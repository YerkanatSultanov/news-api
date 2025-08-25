package router

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "news-api/docs"
	"news-api/internal/http/handlers"
	"news-api/internal/middleware"
	"news-api/pkg/token"
)

func NewRouter(authHandler *handlers.AuthHandler, newsHandler *handlers.NewsHandler, jwtManager *token.JWTManager) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)
	r.Use(mux.CORSMethodMiddleware(r))

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	api.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)

	api.HandleFunc("/news", newsHandler.ListNews).Methods(http.MethodGet)
	api.HandleFunc("/news/{id:[0-9]+}", newsHandler.GetNewsByID).Methods(http.MethodGet)

	secured := api.PathPrefix("").Subrouter()
	secured.Use(middleware.AuthMiddleware(jwtManager))

	secured.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost)

	secured.HandleFunc("/news", newsHandler.CreateNews).Methods(http.MethodPost)
	secured.HandleFunc("/news/{id:[0-9]+}", newsHandler.UpdateNews).Methods(http.MethodPut)
	secured.HandleFunc("/news/{id:[0-9]+}", newsHandler.DeleteNews).Methods(http.MethodDelete)

	return r
}
