package routes

import (
	"net/http"

	"github.com/Romasmi/advanced-url-shortener/internal/application"
	"github.com/Romasmi/advanced-url-shortener/internal/handlers"
	"github.com/Romasmi/advanced-url-shortener/internal/repository"
	"github.com/Romasmi/advanced-url-shortener/internal/services"

	"github.com/gorilla/mux"
)

func RegisterUrlRoutes(router *mux.Router, deps *application.App) {
	urlService := &services.UrlService{UrlRepository: repository.NewUrlRepository(deps.DbConn.DB, deps.RedisConn.Rdb), Config: deps.Config}
	urlHandler := &handlers.UrlHandler{UrlService: urlService}

	router.HandleFunc("/v1/urls", urlHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/s/{shorUrlCode}", urlHandler.RedirectUserToOriginalUrl).Methods(http.MethodGet)
}
