package main

import (
	"github.com/mkrou/crawler/model"
	"github.com/mkrou/crawler/server"
	"log"
)

func main() {
	tasks := model.NewTaskStack()
	go tasks.Crawl()
	log.Fatal(server.Serve(tasks))
}
