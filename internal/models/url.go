package models

import (
	"time"

	uuid "github.com/samborkent/uuidv7"
)

type Url struct {
	ID          uuid.UUID `json:"id"`
	OriginalUrl string    `json:"original_url"`
	Code        string    `json:"code"`
	CreatedAt   time.Time `json:"created_at"`
}

type UrlCreate struct {
	OriginalUrl string `json:"original_url"`
	Code        string `json:"code"`
}
