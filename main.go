package main

import (
	"net/http"
	"log"
	"time"
	"text/template"
)

func main() {
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


}