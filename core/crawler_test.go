package core

import (
	"net/url"
	"strings"
	"testing"

	"github.com/sebastian-nunez/golang-search-engine/utils"
	"golang.org/x/net/html"
)

func TestParsePageBody(t *testing.T) {
	body := strings.NewReader(`
		<html>
			<head>
				<title>Page title 1</title>
				<meta name="description" content="Some page description">
			</head>
			<body>
				<h1>Heading 1</h1>
				<a href="https://internal.com">Internal link</a>
				<a href="https://external.com">External link</a>>
				<a href="/path">Internal relative link</a>
        <h1>Heading 2</h1>
			</body>
		</html>
	`)
	baseURL, _ := url.Parse("https://internal.com")
	wantTitle := "Page title 1"
	wantDesc := "Some page description"
	wantHeadings := "Heading 1, Heading 2"
	wantInternalLinks := []string{"https://internal.com", "https://internal.com/path"}
	wantExternalLinks := []string{"https://external.com"}

	got, err := parsePageBody(body, baseURL)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got.Title != wantTitle {
		t.Errorf("Expected page title '%s', but got '%s'", wantTitle, got.Title)
	}
	if got.Description != wantDesc {
		t.Errorf("Expected page description '%s', but got '%s'", wantDesc, got.Description)
	}
	if got.Headings != wantHeadings {
		t.Errorf("Expected headings '%s', but got '%s'", wantHeadings, got.Headings)
	}
	if !utils.EqualSlices(got.Links.Internal, wantInternalLinks) {
		t.Errorf("Expected internal links '%v', but got '%v'", wantInternalLinks, got.Links.Internal)
	}
	if !utils.EqualSlices(got.Links.External, wantExternalLinks) {
		t.Errorf("Expected external links '%v', but got '%v'", wantExternalLinks, got.Links.External)
	}
}

func TestGetPageLinks(t *testing.T) {
	bodyNode, _ := html.Parse(strings.NewReader(`
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

	got := getPageLinks(bodyNode, baseURL)

	if !utils.EqualSlices(wantInternal, got.Internal) {
		t.Errorf("want internal links %v, but got %v", wantInternal, got.Internal)
	}
	if !utils.EqualSlices(wantExternal, got.External) {
		t.Errorf("want external links %v, but got %v", wantExternal, got.External)
	}
}

func TestGetPageMetadata(t *testing.T) {
	bodyNode, _ := html.Parse(strings.NewReader(`
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

	gotTitle, gotDesc := getPageMetadata(bodyNode)

	if gotTitle != wantTitle {
		t.Errorf("getPageMetadata want title %s but got %s", wantTitle, gotTitle)
	}
	if gotDesc != wantDesc {
		t.Errorf("getPageMetadata want description %s but got %s", wantDesc, gotDesc)
	}
}

func TestGetPageHeadings(t *testing.T) {
	bodyNode, _ := html.Parse(strings.NewReader(`
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

	got := getPageHeadings(bodyNode)

	if got != want {
		t.Errorf("getPageHeadings want %s but got %s", want, got)
	}
}
