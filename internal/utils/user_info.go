package utils

import (
	"net"
	"net/http"
	"strings"

	"github.com/Romasmi/advanced-url-shortener/internal/models"
)

func ExtractUserInfo(r *http.Request) models.UserInfo {
	return models.UserInfo{
		UserAgent: r.UserAgent(),
		IPAddress: getClientIP(r),
		Referrer:  r.Referer(),
	}
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
