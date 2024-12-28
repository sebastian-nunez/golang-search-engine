package core

import (
	"net/url"
	"strings"
	"testing"

	"github.com/sebastian-nunez/golang-search-engine/utils"
	"golang.org/x/net/html"
)

func TestParseBody(t *testing.T) {
	// TODO: fill in
}

func TestGetLinks(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(`
		<html>
			<body>
				<a href="https://internal.com">Internal link</a>
        <a href="https://internal.com/path">Internal link with path</a>
				<a href="https://external.com">External link</a>
        <a href="https://app.internal.com">External link due to subdomain</a>
				<a href="/about">Internal relative link</a>
				<a href="#section">Anchor link</a>
				<a href="mailto:info@internal.com">Mail link</a>
				<a href="tel:+1234567890">Telephone link</a>
				<a href="javascript:void(0)">JavaScript link</a>
				<a href="document.pdf">PDF link</a>
				<a href="document.md">MD link</a>
			</body>
		</html>
	`))
	baseURL, _ := url.Parse("https://internal.com")
	wantInternal := []string{"https://internal.com", "https://internal.com/path", "https://internal.com/about"}
	wantExternal := []string{"https://external.com", "https://app.internal.com"}

	got := getLinks(doc, baseURL)

	if !utils.EqualSlices(wantInternal, got.Internal) {
		t.Errorf("want internal links %v, but got %v", wantInternal, got.Internal)
	}
	if !utils.EqualSlices(wantExternal, got.External) {
		t.Errorf("want external links %v, but got %v", wantExternal, got.External)
	}
}

func TestGetPageMetadata(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(`
		<html>
      <head>
				<title>Some title</title>
				<meta name="description" content="Some page description">
			</head>
			<body>
				<h1>Primary heading</h1>
			</body>
		</html>
	`))
	wantTitle := "Some title"
	wantDesc := "Some page description"

	gotTitle, gotDesc := getPageMetadata(doc)

	if gotTitle != wantTitle {
		t.Errorf("getPageMetadata want title %s but got %s", wantTitle, gotTitle)
	}
	if gotDesc != wantDesc {
		t.Errorf("getPageMetadata want description %s but got %s", wantDesc, gotDesc)
	}
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
