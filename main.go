package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"github.com/arbarlow/sitemapper/scraper"
)

var lock = &sync.Mutex{}

var numberOfWorkers = 3
var workerState = make([]bool, numberOfWorkers)

var rootURL *url.URL
var scraped = map[string]bool{}
var results = map[string]*scraper.Page{}
var errors = map[string]error{}

var queue = make(chan string, 100)
var finished = make(chan bool)

func main() {
	baseURL := os.Getenv("URL")
	if baseURL == "" {
		log.Fatal("URL ENV variable is required")
	}

	url, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	rootURL = url
	log.Printf("Scraping URL %v", url)

	for w := 0; w < numberOfWorkers; w++ {
		go worker(w)
	}

	queue <- url.String()
	<-finished

	f, err := os.Create("sitemap.json")
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.Encode(results)
	log.Printf("Scraped %d URLS, %d Results, %d Errors -- Output: sitemap.json",
		len(scraped), len(results), len(errors))
}

func worker(n int) {
	for u := range queue {
		workerState[n] = true

		// Parse URL and remove trailing slash for uniqueness check
		loc, _ := url.Parse(u)
		path := path.Clean(loc.Path)

		// Lock reading and writing so a worker doesn't pick up the same URL
		lock.Lock()
		_, done := scraped[path]
		lock.Unlock()

		if !done {
			lock.Lock()
			scraped[path] = true
			lock.Unlock()

			page, err := scraper.CrawlPage(rootURL, u)
			if err != nil {
				errors[u] = err
			}

			if page != nil && err == nil {
				results[u] = page

				for _, v := range page.Links {
					go func(link string) {
						queue <- link
					}(v)
				}
			}
		}

		workerState[n] = false
		checkQueueState()
	}
}

// If no workers are working and there are no more urls in the chan, then quit
// Due to channel scheduling randomness we delay this check by one second
func checkQueueState() {
	t := time.NewTimer(time.Second)
	go func() {
		<-t.C
		working := false
		for w := 0; w < numberOfWorkers; w++ {
			if workerState[w] {
				working = true
			}
		}

		if !working && len(queue) == 0 {
			finished <- true
		}
	}()
}
