package crawler

import (
	"../../models"
	"net/http"
	"strings"

	// log "github.com/llimllib/loglevel"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type CrawledPage struct {
	page  models.Page
	links map[string]struct{}
}

type HttpError struct {
	original string
}

func (self HttpError) Error() string {
	return self.original
}

func htmlReader(url string, resp *http.Response) (page models.Page, links map[string]struct{}) {

	links = make(map[string]struct{})
	page.Url = url
	parsedHTML := html.NewTokenizer(resp.Body)

	var start *html.Token
	var text string

	for {
		// Get next token
		_ = parsedHTML.Next()
		token := parsedHTML.Token()

		if token.Type == html.ErrorToken {
			break
		}

		validElement := token.DataAtom == atom.A || token.DataAtom == atom.P || token.DataAtom == atom.Title

		// Save start of valid element
		if validElement && token.Type == html.StartTagToken {
			start = &token
		}

		// Save text of valid element
		if start != nil && token.Type == html.TextToken {
			text = token.Data
		}

		// Manage end of valid element
		if validElement && token.Type == html.EndTagToken {
			if start == nil {
				continue
			}

			switch token.DataAtom {
			case atom.Title:
				page.Title = text
			case atom.P:
				if len(text) > 0 {
					page.Strings = append(page.Strings, text)

				}
			case atom.A:
				for i := range start.Attr {
					if start.Attr[i].Key == "href" {
						link := strings.TrimSpace(start.Attr[i].Val)
						if len(link) > 0 {
							links[link] = struct{}{}
						}
					}
				}
			}

			start = nil
			text = ""
		}

	}
	return
}

// CrawlPage craws page and returns the page relevant information
func CrawlPage(url string) CrawledPage {
	html, err := http.Get(url)

	if err != nil {
		return CrawledPage{}
	}

	page, links := htmlReader(url, html)
	return CrawledPage{page: page, links: links}
}
