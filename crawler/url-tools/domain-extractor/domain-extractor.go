package urltools

import (
	"errors"
	"net/url"
	"strings"
)

func ExtractDomain(urlString string) (string, error) {

	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	if u.Host == "" {
		return "", errors.New("empty host")
	}

	parts := strings.Split(u.Host, ".")
	if len(parts) < 2 {
		return "", errors.New("invalid host")
	}

	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	return domain, nil
}
