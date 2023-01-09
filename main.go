package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/gocolly/colly/v2"
)

func matchAny(urlPath string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, urlPath)
		if matched {
			return true
		} else if err != nil {
			fmt.Println(err)
		}
	}
	return false
}

type flagList []string

func (list *flagList) String() string {
	return fmt.Sprint(*list)
}
func (i *flagList) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	type arrayFlags []string

	var excludePatterns flagList
	var includePatterns flagList

	urlString := flag.String("url", "REQUIRED", "the url to start crawling")
	flag.Var(&excludePatterns, "exclude", "list of regex patterns of url to exclude")
	flag.Var(&includePatterns, "include", "list of regex patterns. This can be used to include a subset of urls, that were excluded via a broad `exclude` pattern")
	flag.Parse()

	startUrl, urlParseError := url.Parse(*urlString)
	if *urlString == "REQUIRED" || urlParseError != nil {
		fmt.Println("invalud startUrl provided (", startUrl, ")")
		exitCode = 2
		return
	}

	fmt.Println(excludePatterns)

	collector := colly.NewCollector(colly.Async())
	collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 8})

	// Find and visit all links on pages with same host as startUrl
	collector.OnHTML("a[href]", func(element *colly.HTMLElement) {
		url := element.Request.URL

		if url.Host == startUrl.Host {
			element.Request.Ctx.Put(element.Attr("href"), url.String())
			element.Request.Visit(element.Attr("href"))
		}
	})

	collector.OnRequest(func(request *colly.Request) {
		if request.URL.Scheme != "https" && request.URL.Scheme != "http" {
			request.Abort()
		}

		urlPath := request.URL.Path
		matchedExclude := matchAny(urlPath, excludePatterns)
		matchedInclude := matchAny(urlPath, includePatterns)

		if matchedExclude && !matchedInclude {
			request.Abort()
		}
	})

	collector.OnError(func(response *colly.Response, err error) {
		if response.StatusCode == 503 || response.StatusCode == 999 || response.StatusCode == 0 {
			// ignore 503 and 999 and 0 status code to avoid flaky errors
			return
		}

		exitCode = 1
		fmt.Println(
			"Error Visiting:\n",
			response.Request.URL,
			"\n",
			err,
			response.StatusCode,
			"\n Found on:",
			response.Ctx.Get(response.Request.URL.String()),
			"\n ",
		)
	})

	collector.Visit(startUrl.String())
	collector.Wait()
}
