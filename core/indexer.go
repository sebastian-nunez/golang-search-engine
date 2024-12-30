package core

import (
	"strings"

	"github.com/sebastian-nunez/golang-search-engine/db"
)

// InvertedIndex is an in-memory inverted index. It maps tokens to a set of page IDs (stored in the database).
type InvertedIndex map[string]map[string]struct{}

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
