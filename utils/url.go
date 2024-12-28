package utils

import "net/url"

func IsSameHost(absoluteURL string, baseURL string) bool {
	absURL, err := url.Parse(absoluteURL)
	if err != nil {
		return false
	}

	baseURLParsed, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	return absURL.Host == baseURLParsed.Host
}
