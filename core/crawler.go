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
	ParsedBody   ParsedBody
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
	baseURL, _ := url.Parse(inputUrl) // Ignoring error since GET request will fail given invalid input URL
	if err != nil || res == nil {
		log.Infof("Something went wrong while crawling '%s': %s", inputUrl, err)
		return CrawlData{URL: inputUrl, Success: false, ResponseCode: 0, ParsedBody: ParsedBody{}}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Infof("Received HTTP status code '%d' while crawling '%s'", res.StatusCode, inputUrl)
		return CrawlData{URL: inputUrl, Success: false, ResponseCode: res.StatusCode, ParsedBody: ParsedBody{}}
	}

	contentType := res.Header.Get(fiber.HeaderContentType)
	if !strings.HasPrefix(contentType, utils.ContentTypeHTML) {
		log.Infof("Received content type of '%s' and expected HTML while crawling '%s'", contentType, inputUrl)
		// Success is set to false since it could be a temporary issue and we can still retry in the future.
		return CrawlData{URL: inputUrl, Success: false, ResponseCode: res.StatusCode, ParsedBody: ParsedBody{}}
	}

	data, err := parseBody(res.Body, baseURL)
	if err != nil {
		log.Infof("Something went wrong getting data from html body for URL '%s': %s", inputUrl, err)
		return CrawlData{URL: inputUrl, Success: false, ResponseCode: res.StatusCode, ParsedBody: ParsedBody{}}
	}

	return CrawlData{URL: inputUrl, Success: true, ResponseCode: res.StatusCode, ParsedBody: data}
}

func parseBody(body io.Reader, baseURL *url.URL) (ParsedBody, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return ParsedBody{}, fmt.Errorf("unable to parse response body: %s", err)
	}

	start := time.Now()

	title, desc := getPageMetadata(doc)
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
	links := Links{}
	if node == nil {
		return links
	}

	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					url, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}
					urlStr := url.String()

					// Check if urlStr is a:
					//  1) Hashtag/anchor
					//  2) Mail link
					//  3) Telephone link
					//  4) JavaScript link
					//  5) PDF or MD file
					if strings.HasPrefix(urlStr, "#") || strings.HasPrefix(urlStr, "mail") || strings.HasPrefix(urlStr, "tel") ||
						strings.HasPrefix(urlStr, "javascript") || strings.HasSuffix(urlStr, ".pdf") || strings.HasSuffix(urlStr, ".md") {
						continue
					}

					if url.IsAbs() { // "valid", full URL (could be internal or external)
						if utils.IsSameHost(url.String(), baseURL.String()) {
							links.Internal = append(links.Internal, url.String())
						} else {
							links.External = append(links.External, url.String())
						}
					} else { // treat as internal URL
						// ResolveReference will help construct a full, absolute URL:
						// baseURL: https://www.example.com
						// url: /about
						// -> rel: https://www.example.com/about
						rel := baseURL.ResolveReference(url)
						links.Internal = append(links.Internal, rel.String())
					}

					break
				}
			}
		}

		child := n.FirstChild
		for child != nil {
			findLinks(child)
			child = child.NextSibling
		}
	}

	findLinks(node)
	return links
}

// Returns (title, description)
func getPageMetadata(node *html.Node) (string, string) {
	if node == nil {
		return "", ""
	}

	title := ""
	description := ""
	// Recursively searching for `meta` tags in the HTML tree and extracts their `name` and `content` attributes.
	var findMetadata func(*html.Node)

	findMetadata = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" && n.FirstChild != nil {
				title = n.FirstChild.Data
			} else if n.Data == "meta" {
				var attrName, attrContent string

				for _, attr := range n.Attr {
					if attr.Key == "name" {
						attrName = attr.Val
					} else if attr.Key == "content" {
						attrContent = attr.Val
					}
				}

				if attrName == "description" {
					description = attrContent
				}
			}
		}

		child := n.FirstChild
		for child != nil {
			findMetadata(child)
			child = child.NextSibling
		}
	}

	findMetadata(node)
	return title, description
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
