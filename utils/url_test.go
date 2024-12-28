package utils

import "testing"

func TestIsSameHost(t *testing.T) {
	testCases := []struct {
		absoluteURL string
		baseURL     string
		want        bool
	}{
		{absoluteURL: "", baseURL: "", want: true},
		{absoluteURL: "", baseURL: "https://google.com", want: false},
		{absoluteURL: "bad URL", baseURL: "https://google.com", want: false},
		{absoluteURL: "https://google.com", baseURL: "bad URL", want: false},
		{absoluteURL: "https://google.com/search?q=hey", baseURL: "https://google.com", want: true},
		{absoluteURL: "https://google.com/path", baseURL: "https://www.google.com", want: false},
		{absoluteURL: "https://google.com", baseURL: "http://google.com", want: true},
		{absoluteURL: "https://google.com", baseURL: "http://google.org", want: false},
	}

	for _, tc := range testCases {
		got := IsSameHost(tc.absoluteURL, tc.baseURL)

		if got != tc.want {
			t.Errorf("For absoluteURL %s and baseURL %s, want %v, but got %v",
				tc.absoluteURL, tc.baseURL, tc.want, got)
		}
	}
}
