package models

import (
	"encoding/json"
	"time"

	uuid "github.com/samborkent/uuidv7"
)

type Url struct {
	ID          uuid.UUID `json:"id,string"`
	OriginalUrl string    `json:"originalUrl"`
	Code        string    `json:"code"`
	CreatedAt   time.Time `json:"createdAt"`
}

type UrlCreate struct {
	OriginalUrl string `json:"originalUrl"`
	Code        string `json:"code"`
}

func (u Url) MarshalJSON() ([]byte, error) {
	type Alias Url
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    u.ID.String(),
		Alias: (*Alias)(&u),
	})
}
