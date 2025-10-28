package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Romasmi/advanced-url-shortener/internal/config"
	"github.com/Romasmi/advanced-url-shortener/internal/kafka"
	"github.com/Romasmi/advanced-url-shortener/internal/models"
	"github.com/Romasmi/advanced-url-shortener/internal/repository"
	uuid "github.com/samborkent/uuidv7"

	"github.com/xyproto/randomstring"
)

type UrlService struct {
	UrlRepository *repository.UrlRepository
	Config        *config.App
	Kafka         *kafka.KafkaConnection
}

func (s *UrlService) Create(ctx context.Context, url string) (*models.Url, error) {
	code, err := s.generateUniqueCode(ctx, s.Config.MaxAttempts)
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

func (s *UrlService) GetByCode(ctx context.Context, code string) (*models.Url, error) {
	return s.UrlRepository.GetByCode(ctx, code)
}

func (s *UrlService) GetRedirectUrl(ctx context.Context, code string, userInfo models.UserInfo) (string, error) {
	url, err := s.GetByCode(ctx, code)
	if err != nil || url == nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", err
		}
		fmt.Printf("unexpected error while retrieveing original URL: %v\n", err)
		return "", err
	}
	go func() {
		clickEvent := models.ClickEvent{
			EventID:   "",
			Domain:    url.OriginalUrl,
			Code:      code,
			Timestamp: "",
			UserAgent: userInfo.UserAgent,
			IPAddress: userInfo.IPAddress,
			Referrer:  userInfo.Referrer,
		}
		clickEventJson, err := json.Marshal(clickEvent)
		if err != nil {
			fmt.Printf("error while converting event data to JSON: %v", clickEventJson)
			return
		}

		err = s.Kafka.Produce("url.clicks", []byte(code), clickEventJson)
		if err != nil {
			fmt.Println(err)
		}
	}()

	return url.OriginalUrl, nil
}

func (s *UrlService) generateUniqueCode(ctx context.Context, maxAttempts uint) (string, error) {
	for i := 0; i < int(maxAttempts); i++ {
		randomString := randomstring.CookieFriendlyString(int(s.Config.UrlCodeLength))
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

func (s *UrlService) GetUrlShortUrl(code string) string {
	domain := strings.TrimSuffix(s.Config.PublicDomain, "/")
	return domain + "/s/" + code
}
