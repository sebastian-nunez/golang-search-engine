package utils

import "testing"

func TestIsSameHost(t *testing.T) {
	testCases := []struct {
		name        string
		absoluteURL string
		baseURL     string
		want        bool
	}{
		{name: "empty urls", absoluteURL: "", baseURL: "", want: true},
		{name: "empty absolute url", absoluteURL: "", baseURL: "https://google.com", want: false},
		{name: "invalid absolute url", absoluteURL: "bad URL", baseURL: "https://google.com", want: false},
		{name: "invalid base url", absoluteURL: "https://google.com", baseURL: "bad URL", want: false},
		{name: "same host with path", absoluteURL: "https://google.com/search?q=hey", baseURL: "https://google.com", want: true},
		{name: "different subdomain", absoluteURL: "https://google.com/path", baseURL: "https://www.google.com", want: false},
		{name: "same host different protocol", absoluteURL: "https://google.com", baseURL: "http://google.com", want: true},
		{name: "different host", absoluteURL: "https://google.com", baseURL: "http://google.org", want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsSameHost(tc.absoluteURL, tc.baseURL)

			if got != tc.want {
				t.Errorf("For absoluteURL %s and baseURL %s, want %v, but got %v",
					tc.absoluteURL, tc.baseURL, tc.want, got)
			}
		})
	}
}
