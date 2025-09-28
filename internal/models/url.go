package models

import "time"

type Url struct {
	ID          string    `json:"id"`
	OriginalUrl string    `json:"original_url"`
	Code        string    `json:"code"`
	CreatedAt   time.Time `json:"created_at"`
}

type UrlCreate struct {
	OriginalUrl string `json:"original_url"`
	Code        string `json:"code"`
}
