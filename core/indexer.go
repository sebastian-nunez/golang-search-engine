package core

import (
	"fmt"

	"github.com/sebastian-nunez/golang-search-engine/db"
)

// Index is an in-memory inverted index. It maps tokens to URL IDs (stored in the database).
type Index map[string][]string

func (idx Index) Add(pages []db.CrawledPage) {
	for _, page := range pages {
		inputStr := fmt.Sprintf("%s %s %s %s", page.URL, page.Title, page.Description, page.Headings)
		tokens := analyze(inputStr)

		for _, token := range tokens {
			urlIDs := idx[token]
			if urlIDs != nil && urlIDs[len(urlIDs)-1] == page.ID {

			}
		}
	}
}
