package routes

import (
	"net/http"
	"shorturl/internal/handlers"
	"shorturl/internal/repository"
	"shorturl/internal/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUrlRoutes(router *mux.Router, db *pgxpool.Pool) {
	urlHandler := &handlers.UrlHandler{UrlService: &services.UrlService{UrlRepository: repository.NewUrlRepository(db)}}

	router.HandleFunc("/v1/urls", urlHandler.Create).Methods(http.MethodPost)
}
