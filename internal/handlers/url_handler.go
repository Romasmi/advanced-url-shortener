package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Romasmi/advanced-url-shortener/internal/repository"
	"github.com/Romasmi/advanced-url-shortener/internal/services"
	"github.com/gorilla/mux"
)

type UrlHandler struct {
	UrlService *services.UrlService
}

type CreateUrlRequest struct {
	Url string `json:"url"`
}

type CreateUrlResponse struct {
	ShortUrl string `json:"shortUrl"`
}

func (h *UrlHandler) Create(w http.ResponseWriter, r *http.Request) {
	var requestPayload CreateUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Original URL is not processed intentionally because it might have some important parameters in query string etc
	newUrl, err := h.UrlService.Create(r.Context(), requestPayload.Url)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			// if it's a duplicate then return existing link
			newUrl, err = h.UrlService.GetByOriginalUrl(r.Context(), requestPayload.Url)
		}
		// TODO log error
		if err != nil {
			fmt.Printf("Error while url creation: %v\n", err)
			http.Error(w, "internal error - try again", http.StatusInternalServerError)
		}
	}

	response := &CreateUrlResponse{
		ShortUrl: h.UrlService.GetUrlShortUrl(newUrl.Code),
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// TODO log error
		fmt.Printf("error while encoding responce %v", err)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UrlHandler) RedirectUserToOriginalUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["shorUrlCode"]
	url, err := h.UrlService.GetRedirectUrl(r.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
