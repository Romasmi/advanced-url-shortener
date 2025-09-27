package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"shorturl/internal/models"
	"shorturl/internal/repository"
	"shorturl/internal/services"
)

type UrlHandler struct {
	UrlService *services.UrlService
}

func (h *UrlHandler) Create(w http.ResponseWriter, r *http.Request) {
	var url models.UrlCreate
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newUrl, err := h.UrlService.Create(r.Context(), &url)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			http.Error(w, "url already exists", http.StatusBadRequest)
			return
		}
		/// TODO log error
		http.Error(w, "internal error - try again", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(newUrl)
	if err != nil {
		// TODO log error
		fmt.Printf("error while encoding responce %v", err)
	}
	w.WriteHeader(http.StatusCreated)
}
