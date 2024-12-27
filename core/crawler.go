package core

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sebastian-nunez/golang-search-engine/utils"
	"golang.org/x/net/html"
)

type CrawlData struct {
	URL          string
	Success      bool
	ResponseCode int
	Body         ParsedBody
}

type ParsedBody struct {
	CrawlTime time.Duration
	// Concatenated string of all <h1> tags
	PageTitle       string
	PageDescription string
	Headings        string
	Links           Links
}

type Links struct {
	// Links within the same domain
	Internal []string
	// Links to websites outside the domain
	External []string
}

func runCrawl(inputUrl string) CrawlData {
	res, err := http.Get(inputUrl)
	baseURL, _ := url.Parse(inputUrl)
	if err != nil || res == nil {
		log.Infof("Something went wrong while crawling '%s': %s", inputUrl, err)
		return CrawlData{URL: inputUrl, Success: false, ResponseCode: 0, Body: ParsedBody{}}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Infof("Received HTTP status code '%d' while crawling '%s'", res.StatusCode, inputUrl)
		return CrawlData{URL: inputUrl, Success: false, ResponseCode: res.StatusCode, Body: ParsedBody{}}
	}

	contentType := res.Header.Get(fiber.HeaderContentType)
	if !strings.HasPrefix(contentType, utils.ContentTypeHTML) {
		log.Infof("Received content type of '%s' and expected HTML while crawling '%s'", contentType, inputUrl)
		// Success is set to false since it could be a temporary issue and we can still retry in the future.
		return CrawlData{URL: inputUrl, Success: false, ResponseCode: res.StatusCode, Body: ParsedBody{}}
	}

	data, err := parseBody(res.Body, baseURL)
	if err != nil {
		log.Infof("Something went wrong getting data from html body for URL '%s': %s", inputUrl, err)
		return CrawlData{URL: inputUrl, Success: false, ResponseCode: res.StatusCode, Body: ParsedBody{}}
	}

	return CrawlData{URL: inputUrl, Success: true, ResponseCode: res.StatusCode, Body: data}
}

func parseBody(body io.ReadCloser, baseURL *url.URL) (ParsedBody, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return ParsedBody{}, fmt.Errorf("unable to parse response body: %s", err)
	}

	start := time.Now()

	title, desc := getPageData(doc)
	headings := getPageHeadings(doc)
	links := getLinks(doc, baseURL)

	end := time.Now()

	return ParsedBody{
		CrawlTime:       end.Sub(start),
		PageTitle:       title,
		PageDescription: desc,
		Headings:        headings,
		Links:           links,
	}, nil
}

func getLinks(node *html.Node, baseURL *url.URL) Links {
	// TODO: fill in
	return Links{}
}

// Returns (title, description)
func getPageData(node *html.Node) (string, string) {
	// TODO: fill in
	return "", ""
}

func getPageHeadings(node *html.Node) string {
	if node == nil {
		return ""
	}

	var headings strings.Builder
	// Recursively explores the HTML tree looking for <h1> tags
	var findH1 func(*html.Node)

	findH1 = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			// Assume the h1 tag will have the text content as its first child
			if n.FirstChild != nil {
				innerText := n.FirstChild.Data

				headings.WriteString(innerText)
				headings.WriteString(", ")
			}
		}

		child := n.FirstChild
		for child != nil {
			findH1(child)
			child = child.NextSibling
		}
	}

	findH1(node)
	// Remove the last comma and space from the concatenated string
	return strings.TrimSuffix(headings.String(), ", ")
}
