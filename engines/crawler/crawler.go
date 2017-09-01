package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"../../models"

	log "github.com/llimllib/loglevel"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type CrawledPage struct { 
	page models.Page 
	links map[string]struct{}
}

type HttpError struct {
	original string
}

func (self HttpError) Error() string {
	return self.original
}

func htmlReader(url string, resp *http.Response) (page models.Page, links map[string]struct{}){
	links = make(map[string]struct{})
	page.Url = url
	parsedHTML := html.NewTokenizer(resp.Body)

	var start *html.Token
	var text string

	for {
		_ = parsedHTML.Next()
		token := parsedHTML.Token()

		if token.Type == html.ErrorToken {
			break
		}

		if start != nil && token.Type == html.TextToken {
			text = fmt.Sprintf("%s%s", text, token.Data)
			// fmt.Println("Token text - %s", text)
			// fmt.Println("Token data - %s", token.Data)
		}

		// Manage title of the page
		if token.DataAtom == atom.Title {
			switch token.Type {
				case html.StartTagToken:
					start = &token

				case html.EndTagToken:
					if start == nil {
						// log.Warnf("Title End found without Start: %s", text)
						continue
					}

					page.Title = text

					start = nil
					text = ""
				}
		}

		if token.DataAtom == atom.A || token.DataAtom == atom.P {
			switch token.Type {
			case html.StartTagToken:
				if token.DataAtom == atom.A && len(token.Attr) > 0 || token.DataAtom == atom.P{
					start = &token
				}

			case html.EndTagToken:
				if start == nil {
					log.Warnf("Link End found without Start: %s", text)
					continue
				}

				if(token.DataAtom == atom.P && len(text) > 0){
					page.Strings = append(page.Strings, text)

				} else if (token.DataAtom == atom.A) {
					for i := range start.Attr {
						if start.Attr[i].Key == "href" {
							link := strings.TrimSpace(start.Attr[i].Val)
							if(len(link) > 0){
								links[link] = struct{}{}
							}
						}
					}
				}

				start = nil
				text = ""
			}
		}
	}
	return
}

func CrawlPage(url string) CrawledPage {
	html, err := getHTML(url)
	
	if err != nil {
		log.Error(err)
		return CrawledPage{}
	}

	page, links := htmlReader(url, html)
	return CrawledPage{page: page, links: links}
}

func getHTML(url string) (resp *http.Response, err error) {
	resp, err = http.Get(url)

	// if resp.StatusCode > 299 {
	// 	err = HttpError{fmt.Sprintf("Error (%d): %s", resp.StatusCode, url)}
	// }

	return
}