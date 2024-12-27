package core

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestParseBody(t *testing.T) {
	// TODO: fill in
}

func TestGetLinks(t *testing.T) {
	// TODO: fill in
}

func TestGetPageData(t *testing.T) {
	// TODO: fill in
}

func TestGetPageHeadings(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(`
		<html>
			<body>
				<h1>Primary heading</h1>
				<div>
					<h1>Secondary heading</h1>
				</div>
				<h2>Should not be included (h2 heading)</h2>
				<h1></h1>
			</body>
		</html>
	`))
	want := "Primary heading, Secondary heading"

	got := getPageHeadings(doc)

	if got != want {
		t.Errorf("getPageHeadings want %s but got %s", want, got)
	}
}
