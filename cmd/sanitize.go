package cmd

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var allowedDomains = map[string]struct{}{
	"velog.io":    {},
	"medium.com":  {},
	"tistory.com": {},
}

func IsAllowedDomain(host string) (bool, string) {
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return false, host
	}
	domain := strings.Join(parts[len(parts)-2:], ".")
	_, allowed := allowedDomains[domain]
	return allowed, domain
}

func ValidateAndSanitizeURL(rawURL string) (string, string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", "", fmt.Errorf("invalid URL scheme: %s", parsedURL.Scheme)
	}
	allowed, host := IsAllowedDomain(parsedURL.Host)
	if !allowed {
		return "", "", fmt.Errorf("domain not allowed: %s", parsedURL.Host)
	}
	parsedURL.RawQuery = ""

	return parsedURL.String(), host, nil
}

func ExtractUserID(messageID string) int {
	parts := strings.Split(messageID, "-")
	if len(parts) < 3 {
		return 0
	}
	lastIndex := len(parts) - 1
	n, err := strconv.Atoi(parts[lastIndex])
	if err != nil {
		return 0
	}
	return n
}
