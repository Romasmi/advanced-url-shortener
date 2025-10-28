package models

type ClickEvent struct {
	EventID   string `json:"eventId"`
	Domain    string `json:"domain"`
	Code      string `json:"code"`
	Timestamp string `json:"timestamp"`
	UserAgent string `json:"userAgent"`
	IPAddress string `json:"ipAddress"`
	Referrer  string `json:"referrer"`
}
