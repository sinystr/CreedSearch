package crawler

import (
	// "strings"
	"../../models"
	"context"
	// "strings"
)

type Engine struct {
	crawlingDeph int
}

func logPage(page CrawledPage){
	println("Page crawled");
	println("Page title: ", page.page.Title)
	println("Page address: ", page.page.Url)
	// println("Found Strings: ")
	// for _, pageString := range page.page.Strings {
	// 	println(pageString)
	// }
}

func (e *Engine) CrawlSite(site string) models.Site {
	pagesToBeCrawled := []string{}
	crawledPages := map[string]struct{}{}
	
	pagesToBeCrawled = append(pagesToBeCrawled, site)

	returnPages := []models.Page{}

	for i := 0; i < e.crawlingDeph; i++{
		parsedPages := crawlPages(pagesToBeCrawled)

		for _, page := range pagesToBeCrawled {
			crawledPages[page] = struct{}{}
		}

		pagesToBeCrawled = pagesToBeCrawled[:0]

		for _, page := range parsedPages {
			returnPages = append(returnPages, page.page)
			logPage(page)
			fitlerLinks(site, &crawledPages, &page.links)
			for link := range page.links {
				pagesToBeCrawled = append(pagesToBeCrawled, link)
			}
		}

	}
	
	return models.Site{ Address: site, Pages: returnPages }
}

func fitlerLinks(site string, crawledPages *map[string]struct{}, links *map[string]struct{}){
	for link := range *links {
		_, found := (*crawledPages)[link]
		// if(!strings.Contains(link, "") || found){
		if(found){
			delete(*links, link)
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
	return &Engine{crawlingDeph: 3}
}