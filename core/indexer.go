package core

import (
	"strings"

	"github.com/sebastian-nunez/golang-search-engine/db"
)

// InvertedIndex is a map that stores tokens (words) as keys and maps each token
// to a set of associated page IDs. This structure allows for efficient lookups
// of pages containing specific tokens.
//
// Example:
//
//	index := InvertedIndex{
//	    "hello": {"page1": struct{}{}, "page3": struct{}{}},
//	    "world": {"page2": struct{}{}, "page3": struct{}{}},
//	}
//
// This example shows that the token "hello" appears on pages "page1" and "page3",
// while the token "world" appears on pages "page2" and "page3".
type InvertedIndex map[string]map[string]struct{}

// Add adds crawled pages to the inverted index (in-memory), efficiently handling
// duplicate page entries for the same token.
func (idx InvertedIndex) Add(pages []db.CrawledPage) {
	for _, page := range pages {
		doc := buildDocument(page)
		tokens := createIndexTokens(doc)

		for _, token := range tokens {
			_, ok := idx[token]
			if !ok {
				idx[token] = make(map[string]struct{})
			}

			idx[token][page.ID] = struct{}{}
		}
	}
}

// buildDocument concatenates relevant text fields for indexing.
func buildDocument(page db.CrawledPage) string {
	fields := []string{page.URL, page.Title, page.Description, page.Headings}
	return strings.Join(fields, " ")
}
