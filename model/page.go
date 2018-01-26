package model

import (
	"sync"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"time"
	"net/url"
)

const (
	CodeNotUrl   = -1
	CodeNetError = 0
)

type page struct {
	Title string `json:"title,omitempty"`
	Code  int    `json:"code"`
	Url   string `json:"url"`
}

func (p *page) parse(wg *sync.WaitGroup, sync *Throttler) {
	defer wg.Done()
	if !urlValidate(p.Url) {
		p.Code = CodeNotUrl
		return
	}
	address, _ := url.Parse(p.Url)
	time.Sleep(sync.Delay(address.Hostname()))
	resp, err := http.Get(p.Url)
	if err != nil {
		p.Code = CodeNetError
		return
	}
	p.Code = resp.StatusCode
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		p.Code = CodeNetError
		return
	}
	p.Title = doc.Find("title").Text()
}

func pagesFromUrls(urls []string) []*page {
	var pages []*page
	for _, url := range urls {
		pages = append(pages, &page{
			Url: url,
		})
	}
	return pages
}

//This validation works better than net/url.ParseRequestURI
func urlValidate(url string) bool {
	return regexp.MustCompile(`https?://(-\.)?([^\s/?\.#-]+\.?)+(/[^\s]*)?$`).MatchString(url)
}
