package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Page struct
type Page struct {
	url          string
	title        string
	lastModified string
	pageSize     string
	keywords     []string
	parentURL    []string
	childrenURL  []string
}

// GetURL return url of page
func (page *Page) GetURL() string {
	return page.url
}

// GetTitle return title of the page
func (page *Page) GetTitle() string {
	return page.title
}

// GetLastModified return the last modified date of a page
func (page *Page) GetLastModified() string {
	return page.lastModified
}

// GetSize return size of page
func (page *Page) GetSize() string {
	return page.pageSize
}

// GetKeywords return list of keywords
func (page *Page) GetKeywords() []string {
	return page.keywords
}

// GetChildrenURL return list of keywords
func (page *Page) GetChildrenURL() []string {
	return page.childrenURL
}

// GetParentURL return url of page
func (page *Page) GetParentURL() []string {
	return page.parentURL
}

// ExtractTitle from each url
func (page *Page) ExtractTitle() {
	res, err := http.Get(page.GetURL())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// log.Fatalf("status code error title: %d %s, %s", res.StatusCode, res.Status, page.GetURL())
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	page.title = doc.Find("title").Text()
}

// ExtractLastModified extract the last-modified date from a header of a url, return 0 if not found
func (page *Page) ExtractLastModified() {
	res, err := http.Head(page.GetURL())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// log.Fatalf("status code error title: %d %s, %s", res.StatusCode, res.Status, page.GetURL())
	}
	if res.Header.Get("Date") != "" {
		if res.Header.Get("Last-Modified") != "" {
			page.lastModified = res.Header.Get("Last-Modified")
		} else {
			page.lastModified = res.Header.Get("Date")
		}
	} else {
		page.lastModified = "0"
	}
}

// ExtractSize extract the size of the page if found, -1 if not found
func (page *Page) ExtractSize() {
	res, err := http.Head(page.GetURL())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// log.Fatalf("status code error title: %d %s, %s", res.StatusCode, res.Status, page.GetURL())
	}
	if res.Header.Get("Content-Length") != "" {
		page.pageSize = res.Header.Get("Content-Length")
	} else {
		page.pageSize = "-1"
	}
}

// ExtractWords from each url
func (page *Page) ExtractWords() {
	res, err := http.Get(page.GetURL())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// log.Fatalf("status code error keywords: %d %s, %s", res.StatusCode, res.Status, page.GetURL())
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodytext := doc.Find("body").Text()
	if doc.Find("main").Text() != "" {
		bodytext = doc.Find("main").Text()
	}
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
					// make sure there's no duplicate e.g. http://www.cse.ust.hk//pg
					if string(baseURL[len(baseURL)-1]) == "/" {
						newurl := baseURL[:len(baseURL)-1]
						url = newurl + url
						ch <- url
					} else {
						url = baseURL + url
						ch <- url
					}

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
	for url := range foundUrls {
		page.childrenURL = append(page.childrenURL, url)
	}

	close(chUrls)
}

// MakeChildren given a page, create its children page and append it into the same slice with the parent
func (page *Page) MakeChildren(pages *[]*Page) {
	for _, url := range page.GetChildrenURL() {
		childPage := Page{url, "", "", "", make([]string, 0), make([]string, 0), make([]string, 0)}
		childPage.ExtractTitle()
		childPage.ExtractLastModified()
		childPage.ExtractSize()
		childPage.ExtractWords()
		childPage.ExtractLinks()
		childPage.parentURL = append(childPage.parentURL, page.GetURL())
		*pages = append(*pages, &childPage)
	}
}

// WriteIndexed write the result of extraction into a file.txt
func WriteIndexed(pages *[]*Page) {
	(*pages)[0].ExtractTitle()
	(*pages)[0].ExtractLastModified()
	(*pages)[0].ExtractSize()
	(*pages)[0].ExtractWords()
	(*pages)[0].ExtractLinks()
	f, err := os.Create("spider_result.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	head := "TITLE: " + (*pages)[0].GetTitle() + "\n" + (*pages)[0].GetURL() + "\n" + "DATE: " + (*pages)[0].GetLastModified() + ", " + (*pages)[0].GetSize() + "\n" + strings.Join((*pages)[0].GetKeywords(), " ") + "\n" + strings.Join((*pages)[0].GetChildrenURL(), "\n") + "\n"
	(*pages)[0].MakeChildren(pages)
	for i := 1; i < len(*pages); i++ {
		if (*pages)[i].GetTitle() != "" {
			children := "----------------------------------------------\n" + "TITLE: " + (*pages)[i].GetTitle() + "\n" + (*pages)[i].GetURL() + "\n" + "DATE: " + (*pages)[i].GetLastModified() + ", " + (*pages)[i].GetSize() + "\n" + strings.Join((*pages)[i].GetKeywords(), " ") + "\n" + strings.Join((*pages)[i].GetChildrenURL(), "\n") + "\n"
			head += children
		}
	}

	_, err = f.WriteString(head)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	// fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	const baseURL = "https://www.cse.ust.hk/"
	fmt.Println(time.Now())

	pages := make([]*Page, 0)
	basePage := Page{baseURL, "", "", "", make([]string, 0), nil, make([]string, 0)}
	pages = append(pages, &basePage)
	WriteIndexed(&pages)

	fmt.Println(time.Now())
}

// javascript:alert(document.lastModified)
// https://www.techinasia.com/top-funded-startups-tech-companies-india?ref=subexc-444416

/*
	// res, err := http.Get("https://www.cse.ust.hk/")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// robots, err := ioutil.ReadAll(res.Body)
	// res.Body.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s", robots)
*/
