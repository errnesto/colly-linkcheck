package main

import (
	"fmt"
	"strings"

	// "regexp"

	"github.com/gocolly/colly/v2"
)

func main() {
	collector := colly.NewCollector(
		colly.Async(),
	)

	collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})

	// Find and visit all links on ems.press pages
	collector.OnHTML("a[href]", func(element *colly.HTMLElement) {
		if element.Request.URL.Host == "ems.press" {
			element.Request.Visit(element.Attr("href"))
		}
	})

	collector.OnRequest(func(request *colly.Request) {
		if request.URL.Scheme == "mailto" ||
			strings.HasPrefix(request.URL.Path, "/journals") ||
			strings.HasPrefix(request.URL.Path, "/books") ||
			strings.HasPrefix(request.URL.Path, "/content") {
			request.Abort()
		}
	})

	collector.OnError(func(response *colly.Response, err error) {
		fmt.Println("Error Visiting:", response.Request.URL, err, response.StatusCode)
	})

	collector.Visit("https://ems.press/")
	collector.Wait()
}
