package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Page struct
type Page struct {
	url         string
	title       string
	keywords    []string
	childrenURL []string
}

// GetURL return url of crawler
func (page *Page) GetURL() string {
	return page.url
}

// GetTitle return list of keywords
func (page *Page) GetTitle() string {
	return page.title
}

// GetKeywords return list of keywords
func (page *Page) GetKeywords() []string {
	return page.keywords
}

// GetChildrenURL return list of keywords
func (page *Page) GetChildrenURL() []string {
	return page.childrenURL
}

// ExtractTitle from each url
func (page *Page) ExtractTitle() string {
	res, err := http.Get(page.GetURL())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// rows := make([]string, 0)

	page.title = doc.Find("title").Text()
	// rows = append(rows, title)
	return page.title
}

// ExtractWords from each url
func (page *Page) ExtractWords() {
	res, err := http.Get(page.GetURL())
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

	bodytext := doc.Find("body").Text()
	text := strings.Fields(bodytext)
	for _, keyword := range text {
		page.keywords = append(page.keywords, keyword)
	}

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
			hasProto := strings.Index(url, "http") == 0
			startSlash := strings.Index(url, "/") == 0
			if hasProto || startSlash {
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
func (page *Page) ExtractLinks() {
	foundUrls := make(map[string]bool)

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	go getLinks(page.GetURL(), chUrls, chFinished)

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
	// i := 0
	for url := range foundUrls {
		page.childrenURL = append(page.childrenURL, url)
		// if i == 29 {
		// 	break
		// }
		// i++
	}

	close(chUrls)
}

// ExtractAll extract all words, keywords, title, etc from a given url and its immediate children
func (page *Page) ExtractAll() {
	page.ExtractTitle()
	page.ExtractWords()
	page.ExtractLinks()
}

// WriteIndexed write the result of extraction into a file.txt
func (page *Page) WriteIndexed() {
	page.ExtractAll()
	f, err := os.Create("spider_result.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	head := page.GetTitle() + "\n" + page.GetURL() + "\n" + strings.Join(page.GetKeywords(), " ") + "\n" + strings.Join(page.GetChildrenURL(), "\n") + "\n"

	for _, url := range page.GetChildrenURL() {
		childPage := Page{url, "", make([]string, 0), make([]string, 0)}
		childPage.ExtractAll()
		// fmt.Println(childPage.GetTitle())
		if childPage.GetTitle() != "" {
			children := "---------------------------------------\n" + childPage.GetTitle() + "\n" + childPage.GetURL() + "\n" + strings.Join(childPage.GetKeywords(), " ") + "\n" + strings.Join(childPage.GetChildrenURL(), "\n") + "\n"
			head += children
		}
	}

	l, err := f.WriteString(head)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	const baseURL = "https://www.cse.ust.hk/"
	page := Page{baseURL, "", make([]string, 0), make([]string, 0)}
	page.WriteIndexed()
}
