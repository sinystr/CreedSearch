package main

import (
	"fmt"
	"net/http"
	"log"
	"time"
	"text/template"
	"./engines/crawler"
)

// Creed is search engine written in Go
type Creed struct {
	
	// databaseEngine is going to be used for managing the creed database
	DatabaseEngine DatabaseEngine

	// crawlingEngine is going to be used for crawling sites
	CrawlingEngine CrawlingEngine

	// searchEngine is going to be used for text in site's retrieved content
	SearchEngine SearchEngine

}

func (c *Creed) SetDatabaseEngine(databaseEngine DatabaseEngine) {
	c.DatabaseEngine = databaseEngine
}

func (c *Creed) SetCrawlingEngine(crawlingEngine CrawlingEngine) {
	c.CrawlingEngine = crawlingEngine
}

func (c *Creed) SetSearchEngine(searchEngine SearchEngine) {
	c.SearchEngine = searchEngine
}

func (c *Creed) startServer(port int) {
	mux := http.NewServeMux()
	
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/" {
						t, _ := template.ParseFiles("index.html")
						t.Execute(w, nil)
					return
			}
			if req.URL.Path == "/search" {
				searchSite := req.URL.Query()["site"][0]
				var crawler CrawlingEngine = crawler.DefaultEngine()
				site, err := crawler.CrawlSite(searchSite)

				if(err != nil) {
					t, _ := template.ParseFiles("message.html")
					t.Execute(w, "Сайта не може да бъде достъпен")
					return;	
				}

				t, _ := template.ParseFiles("search_results.html")
				t.Execute(w, site)
				return;	
			}
			http.NotFound(w, req)
	})
	
	s := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux, WriteTimeout: 1 * time.Second}
	log.Printf("Starting server on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
	
	fmt.Println("Starting server...");
}

// DefaultCreed returns the default configuration Creed
func DefaultCreed() *Creed {
	creed := Creed{}
	creed.SetCrawlingEngine(crawler.DefaultEngine())
	return &creed;
}