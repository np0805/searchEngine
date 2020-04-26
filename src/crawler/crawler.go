package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Crawler struct
type Crawler struct {
	url string
}

// GetURL return url of crawler
func (crawler *Crawler) GetURL() string {
	return crawler.url
}

// GetTitle from each url
func (crawler *Crawler) GetTitle() {
	res, err := http.Get(crawler.GetURL())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	rows := make([]string, 0)

	title := doc.Find("title").Text()
	rows = append(rows, title)
	fmt.Println(rows)
}

// getHref get the href attribute from the token
func getHref(token html.Token) (exist bool, href string) {
	for _, x := range token.Attr {
		if x.Key == "href" {
			href = x.Val
			exist = true
		}
	}
	return exist, href
}

// getLinks get the links from the given url
func getLinks(url string, ch chan string, chFinished chan bool) {
	baseURL := url
	resp, err := http.Get(baseURL)

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl:", url)
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function completes

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}

			// Make sure the url begines in http**
			hasProto := strings.Index(url, "#") == 0
			startSlash := strings.Index(url, "/") == 0
			if !hasProto {
				if startSlash {
					url = baseURL + url
					ch <- url
				} else {
					ch <- url
				}
			}
		}
	}
}

// ExtractLinks get the links from the url
func (crawler *Crawler) ExtractLinks() {
	foundUrls := make(map[string]bool)
	// seedUrls := os.Args[1:]

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	// for _, url := range seedUrls {
	// 	go crawl(crawler.GetURL(), chUrls, chFinished)
	// }
	go getLinks(crawler.GetURL(), chUrls, chFinished)

	// Subscribe to both channels
	for c := 0; c < 1; {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		}
	}

	// We're done! Print the results...

	fmt.Println("\nFound", len(foundUrls), "unique urls in : ", crawler.GetURL(), "\n")

	for url := range foundUrls {
		fmt.Println(" - " + url)
	}

	close(chUrls)
}

// ExtractWords extract words from a given link
func (crawler *Crawler) ExtractWords() []string {
	words := make([]string, 0)
	words = append(words, crawler.GetURL())
	return words
}

func main() {
	const baseURL = "https://www.cse.ust.hk/"
	crawler := Crawler{baseURL}
	crawler.GetTitle()
	crawler.ExtractLinks()
}
