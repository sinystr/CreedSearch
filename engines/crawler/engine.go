package crawler

import (
	"../../models"
	"context"
	"net/url"
	"strings"
)

// Engine defines crawling engine for Creed
type Engine struct {
	crawlingDeph int
}

// CrawlSite accepts string value of a website address and returns
// crawled site structure
func (e *Engine) CrawlSite(site string) (models.Site, error) {
	siteURL, err := url.ParseRequestURI(site)

	if err != nil {
		return models.Site{}, err
	}

	pagesToBeCrawled := []string{}
	crawledPages := map[string]struct{}{}

	pagesToBeCrawled = append(pagesToBeCrawled, site)

	returnPages := []models.Page{}

	for i := 0; i < e.crawlingDeph; i++ {
		parsedPages := crawlPages(pagesToBeCrawled)

		for _, page := range pagesToBeCrawled {
			crawledPages[page] = struct{}{}
		}

		pagesToBeCrawled = pagesToBeCrawled[:0]

		for _, page := range parsedPages {
			returnPages = append(returnPages, page.page)

			filterLinks(siteURL, &crawledPages, &page.links)
			for link := range page.links {
				pagesToBeCrawled = append(pagesToBeCrawled, link)
			}
		}

	}

	return models.Site{Address: site, Pages: returnPages}, nil
}

func filterLinks(site *url.URL, crawledPages *map[string]struct{}, links *map[string]struct{}) {
	for link := range *links {
		_, found := (*crawledPages)[link]

		if found {
			delete(*links, link)
		} else {
			if link[0] == '/' {
				delete(*links, link)
				newLink := site.String() + link[1:len(link)]
				(*links)[newLink] = struct{}{}
			} else {
				if !strings.Contains(link, site.Hostname()) {
					delete(*links, link)
				}
			}
		}
	}
}

func crawlPages(pages []string) (crawledPages []CrawledPage) {
	crawledPages = make([]CrawledPage, 0, len(pages))
	ctx, cancel := context.WithCancel(context.Background())
	resultChannel := CrawlPages(ctx, pages)
	defer cancel()

	for crawledPage := range resultChannel {
		crawledPages = append(crawledPages, crawledPage)
	}

	return crawledPages
}

// DefaultEngine return the crawling engine using the default configuration
func DefaultEngine() *Engine {
	return &Engine{crawlingDeph: 2}
}

// EngineWithDepth returns crawing engine with custom depth
func EngineWithDepth(depth int) *Engine {
	return &Engine{crawlingDeph: depth}
}
