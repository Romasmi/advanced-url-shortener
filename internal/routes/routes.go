package routes

import (
	"encoding/json"
	"github.com/Romasmi/advanced-url-shortener/internal/application"
	"net/http"

	"github.com/gorilla/mux"
)

type NotFoundResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(router *mux.Router, deps *application.App) {
	if router == nil {
		panic("router must be initialized before routes registration")
	}

	RegisterUrlRoutes(router, deps)

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	response := &NotFoundResponse{Error: "Route Not found"}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}
