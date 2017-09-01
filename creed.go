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
	databaseEngine *DatabaseEngine

	// crawlingEngine is going to be used for crawling sites
	crawlingEngine *CrawlingEngine

	// searchEngine is going to be used for text in site's retrieved content
	searchEngine *SearchEngine

}

func (c *Creed) setDatabaseEngine(databaseEngine *DatabaseEngine) {
	c.databaseEngine = databaseEngine
}

func (c *Creed) setCrawlingEngine(crawlingEngine *CrawlingEngine) {
	c.crawlingEngine = crawlingEngine
}

func (c *Creed) setSearchEngine(searchEngine *SearchEngine) {
	c.searchEngine = searchEngine
}

func (c *Creed) startServer() {
	var crawler CrawlingEngine = crawler.DefaultEngine()
	crawler.CrawlSite("http://fmi.ruby.bg/")
	return;

	mux := http.NewServeMux()
	
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/" {
						t, _ := template.ParseFiles("index.html")
						t.Execute(w, nil)
					return
			}
			http.NotFound(w, req)
	})
	
	s := &http.Server{Addr: ":8282", Handler: mux, WriteTimeout: 1 * time.Second}
	log.Printf("Starting server on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
	
	fmt.Println("Starting server...");
}

// func StrandartCreed() *Creed {
// 	// return;
// }