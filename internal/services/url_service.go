package services

import (
	"context"
	"shorturl/internal/models"
	"shorturl/internal/repository"
)

type UrlService struct {
	UrlRepository *repository.UrlRepository
}

func (s *UrlService) Create(ctx context.Context, url *models.Url) (*models.Url, error) {
	if url.ID == "" {
		url.ID = "generated guid"
	}
	return s.UrlRepository.Create(ctx, url)
}

func (s *UrlService) GetByOriginalUrl(ctx context.Context, originalUrl string) (*models.Url, error) {
	return s.UrlRepository.GetByOriginalUrl(ctx, originalUrl)
}
