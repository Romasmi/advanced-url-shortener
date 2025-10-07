package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/advanced-url-shortener/internal/config"
	"github.com/Romasmi/advanced-url-shortener/internal/models"
	"github.com/Romasmi/advanced-url-shortener/internal/repository"
	uuid "github.com/samborkent/uuidv7"

	"github.com/xyproto/randomstring"
)

type UrlService struct {
	UrlRepository *repository.UrlRepository
	Config        *config.Config
}

func (s *UrlService) Create(ctx context.Context, url string) (*models.Url, error) {
	code, err := s.generateUniqueCode(ctx, s.Config.App.MaxAttempts)
	if err != nil {
		return nil, err
	}

	return s.UrlRepository.Create(ctx, &models.Url{
		ID:          uuid.New(),
		OriginalUrl: url,
		Code:        code,
	})
}

func (s *UrlService) GetByOriginalUrl(ctx context.Context, originalUrl string) (*models.Url, error) {
	return s.UrlRepository.GetByOriginalUrl(ctx, originalUrl)
}

func (s *UrlService) generateUniqueCode(ctx context.Context, maxAttempts uint) (string, error) {
	for i := 0; i < int(maxAttempts); i++ {
		randomString := randomstring.CookieFriendlyString(int(s.Config.App.UrlCodeLength))
		_, err := s.UrlRepository.GetByCode(ctx, randomString)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return randomString, nil
			}
			return "", err
		}
	}
	return "", fmt.Errorf("failed to generated unique code after %v attempts", maxAttempts)
}
