package scraper

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	URL    string
	Assets []string
	Links  []string
}

func CrawlPage(rootURL *url.URL, loc string) (*Page, error) {
	fmt.Printf("loc = %+v\n", loc)
	page := &Page{URL: loc}

	doc, err := goquery.NewDocument(loc)
	if err != nil {
		return nil, err
	}

	findAssets(page, doc)
	findLinks(rootURL, page, doc)

	return page, nil
}

func findAssets(page *Page, doc *goquery.Document) {
	doc.Find("script, img").Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr("src")
		if !ok {
			return
		}
		page.Assets = append(page.Assets, val)
	})

	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr("href")
		if !ok {
			return
		}
		page.Assets = append(page.Assets, val)
	})
}

func findLinks(rootURL *url.URL, page *Page, doc *goquery.Document) {
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr("href")
		if !ok {
			return
		}

		fullURL, isInternal, err := isInternalURL(rootURL, val)
		if err != nil {
			return
		}

		if isInternal {
			page.Links = append(page.Links, fullURL)
		}
	})
}

func isInternalURL(rootURL *url.URL, loc string) (string, bool, error) {
	hurl, err := url.Parse(loc)
	if err != nil {
		return "", false, err
	}

	if hurl.Host == "" || hurl.Host == rootURL.Host {
		if hurl.Path != "/" {
			// We have already parsed the rootURL, no error should occur
			crawlURL, _ := url.Parse(rootURL.String())
			crawlURL.Path = hurl.Path
			return crawlURL.String(), true, nil
		}
	}

	return loc, false, nil
}
