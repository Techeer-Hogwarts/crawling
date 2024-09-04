package cmd

import (
	"fmt"
	"net/url"
)

var allowedDomains = map[string]struct{}{
	"velog.io":   {},
	"medium.com": {},
}

func IsAllowedDomain(host string) bool {
	_, allowed := allowedDomains[host]
	return allowed
}

func ValidateAndSanitizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("invalid URL scheme: %s", parsedURL.Scheme)
	}

	if !IsAllowedDomain(parsedURL.Host) {
		return "", fmt.Errorf("domain not allowed: %s", parsedURL.Host)
	}
	parsedURL.RawQuery = ""

	return parsedURL.String(), nil
}
