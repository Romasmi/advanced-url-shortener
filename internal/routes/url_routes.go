package routes

import (
	"net/http"

	"github.com/Romasmi/advanced-url-shortener/internal/config"
	"github.com/Romasmi/advanced-url-shortener/internal/handlers"
	"github.com/Romasmi/advanced-url-shortener/internal/repository"
	"github.com/Romasmi/advanced-url-shortener/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/gorilla/mux"
)

func RegisterUrlRoutes(
	router *mux.Router,
	db *pgxpool.Pool,
	redis *redis.Client,
	config *config.Config,
) {
	urlService := &services.UrlService{UrlRepository: repository.NewUrlRepository(db, redis), Config: config}
	urlHandler := &handlers.UrlHandler{UrlService: urlService}

	router.HandleFunc("/v1/urls", urlHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/s/{shorUrlCode}", urlHandler.RedirectUserToOriginalUrl).Methods(http.MethodGet)
}
