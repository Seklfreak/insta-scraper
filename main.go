package main

import (
	"path"
	"strings"

	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

const (
	start = "https://www.instagram.com/nayoung_lim95/"
	dir   = "./result/"
)

var (
	log *zap.Logger
	c   *colly.Collector
)

func main() {
	var err error
	log, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	c = colly.NewCollector()

	c.OnHTML("body > script", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Text, "window._sharedData") {
			crawlPage(e)
			return
		}
		if strings.HasPrefix(e.Text, "window.__additionalDataLoaded") {
			crawlAdditionalData(e)
			return
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Debug("visiting page", zap.Stringer("url", r.URL))
	})

	c.OnResponse(func(r *colly.Response) {
		if strings.Contains(r.Headers.Get("Content-Type"), "image/") ||
			strings.Contains(r.Headers.Get("Content-Type"), "video/") {
			err = r.Save(dir + path.Base(r.Request.URL.Path))
			if err != nil {
				log.Error("failure saving image", zap.Error(err))
			}
			return
		}

		if strings.Contains(r.Headers.Get("Content-Type"), "application/json") {
			crawlProfilePostsGraphQL(r.Request.URL, r.Body)
			return
		}
	})

	err = c.Visit(start)
	if err != nil {
		log.Error("crawler failed", zap.Error(err))
	}
}
