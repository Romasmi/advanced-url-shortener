package services

import (
	"context"
	"shorturl/internal/models"
	"shorturl/internal/repository"
)

type UrlService struct {
	UrlRepository *repository.UrlRepository
}

func (s *UrlService) Create(ctx context.Context, url *models.UrlCreate) (*models.Url, error) {
	return s.UrlRepository.Create(ctx, &models.Url{
		ID:          "some id",
		OriginalUrl: url.OriginalUrl,
		ShortUrl:    url.ShortUrl,
	})
}

func (s *UrlService) GetByOriginalUrl(ctx context.Context, originalUrl string) (*models.Url, error) {
	return s.UrlRepository.GetByOriginalUrl(ctx, originalUrl)
}
