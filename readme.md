# Go Simple Web Mapper

A simple, fairly hacky, fairly quickly made web scraper.

Can find, scripts, stylesheets, images and links in any document and outputs a JSON file.

## Usage
```
go get github.com/arbarlow/sitemapper
URL=https://example.com/ sitemapper
```

When finished the scraper will output an 'sitemap.json' file and by default using 3 concurrent scrapers.

Feel free to play around with `numberOfWorkers` that variable in `main.go`

## Issues
Doesn't parse CSS/JS to calculate dependencies that are loaded via dependencies.

Should be kinder to websites with a rate limiter

On error, no retry. Probably should retry with an exponential backoff

Structures are arrays, not maps. So links and assets are not deduped.

There are no tests, I'm a bad person.
