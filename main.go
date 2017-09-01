package main

import (
	// "strings"
	"time"
)

func main() {

	println("Start crawling: ", time.Now().String())
	
	creedServer := &Creed{}
	creedServer.startServer()

    println("End crawling: ", time.Now().String())
}