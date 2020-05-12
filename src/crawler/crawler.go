package crawler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Page struct
type Page struct {
	URL          string
	Title        string
	LastModified string
	PageSize     string
	pageRank     float64
	Keywords     []string
	ParentURL    []string
	ChildrenURL  []string
}

// GetURL return url of page
func (page *Page) GetURL() string {
	return page.URL
}

// GetTitle return title of the page
func (page *Page) GetTitle() string {
	return page.Title
}

// GetLastModified return the last modified date of a page
func (page *Page) GetLastModified() string {
	return page.LastModified
}

// GetSize return size of page
func (page *Page) GetSize() string {
	return page.PageSize
}

// GetKeywords return list of keywords
func (page *Page) GetKeywords() []string {
	return page.Keywords
}

// GetChildrenURL return url of its children
func (page *Page) GetChildrenURL() []string {
	return page.ChildrenURL
}

// GetParentURL return url of its parents
func (page *Page) GetParentURL() []string {
	return page.ParentURL
}

// GetPageRank return page rank
func (page *Page) GetPageRank() float64 {
	return page.pageRank
}

// SetRank sets the calling page with the given rank
func (page *Page) SetRank(rank float64) {
	page.pageRank = rank
}

// ExtractTitle extract the title from a given page
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
	page.Title = doc.Find("title").Text()
}

// ExtractLastModified extract the last-modified date from a header of a given page, return 0 if not found
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
			page.LastModified = res.Header.Get("Last-Modified")
		} else {
			page.LastModified = res.Header.Get("Date")
		}
	} else {
		page.LastModified = "0"
	}
}

// ExtractSize extract the size of the page if found, if not found then value will be size of characters in the page
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
		page.PageSize = res.Header.Get("Content-Length")
	} else {
		count := 0
		for _, words := range page.Keywords {
			count += utf8.RuneCountInString(words)
		}
		page.PageSize = strconv.Itoa(count)
	}
}

// ExtractWords extract keywords from a given page
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
		page.Keywords = append(page.Keywords, keyword)
	}
}

// getHref get the href attribute from the token from a url
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
	resp, err := http.Get(url)

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

			// Make sure the url begines in https://www.cse.ust.hk/ as we only index from those pages
			hasProto := strings.Index(url, "https://www.cse.ust.hk/") == 0
			startSlash := strings.Index(url, "/") == 0
			cseHK := "https://www.cse.ust.hk"
			if hasProto || startSlash {
				if startSlash {
					// make sure there's no duplicate "/" e.g. http://www.cse.ust.hk//pg
					// if string(cseHK[len(baseURL)-1]) == "/" {
					// 	newurl := baseURL[:len(baseURL)-1]
					// 	url = newurl + url
					// 	ch <- url
					// } else {
					url = cseHK + url
					ch <- url
					// }

				} else {
					ch <- url
				}
			} else {
				continue
			}
		}
	}
}

// ExtractLinks get the links from the given page
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
		page.ChildrenURL = append(page.ChildrenURL, url)
	}

	close(chUrls)
}

// MakeLessChildren given a page, create 30 of its children page from the page's childrenURL and map it to the given map
func (page *Page) MakeLessChildren(pages *map[string]*Page) {
	for i, url := range page.GetChildrenURL() {
		// check if it is from cse.ust.hk
		// if strings.Index(url, "https://www.cse.ust.hk/") == 0 {
		childPage, ok := (*pages)[url]
		if !ok {
			childPage := Page{url, "", "", "", 1, make([]string, 0), make([]string, 0), make([]string, 0)}
			childPage.ExtractTitle()
			if childPage.GetTitle() == "" {
				continue
			}
			childPage.ExtractLastModified()
			childPage.ExtractWords()
			childPage.ExtractSize()
			childPage.ExtractLinks()
			childPage.ParentURL = append(childPage.ParentURL, page.GetURL())
			(*pages)[childPage.URL] = &childPage
		} else {
			// Check for circular dependency
			if childPage.GetURL() == page.GetURL() {
				// fmt.Println("Yaha kamu ketauan")
				continue
			} else {
				// fmt.Println("eits sudah pernah bro")
				childPage.ParentURL = append(childPage.ParentURL, page.GetURL())
			}
		}
		if i == 5 {
			break
		}
	}
}

// MakeChildren given a page, create its children page from the page's childrenURL and map it to the given map
func (page *Page) MakeChildren(pages *map[string]*Page) {
	for _, url := range page.GetChildrenURL() {
		// check if it is from cse.ust.hk
		// if strings.Index(url, "https://www.cse.ust.hk/") == 0 {
		childPage, ok := (*pages)[url]
		if !ok {
			childPage := Page{url, "", "", "", 1, make([]string, 0), make([]string, 0), make([]string, 0)}
			childPage.ExtractTitle()
			if childPage.GetTitle() == "" {
				continue
			}
			childPage.ExtractLastModified()
			childPage.ExtractWords()
			childPage.ExtractSize()
			childPage.ExtractLinks()
			childPage.ParentURL = append(childPage.ParentURL, page.GetURL())
			(*pages)[childPage.URL] = &childPage
			// Computer recursively
			childPage.MakeChildren(pages)
		} else {
			// Check for circular dependency
			if childPage.GetURL() == page.GetURL() {
				// fmt.Println("Yaha kamu ketauan")
				continue
			} else {
				// fmt.Println("eits sudah pernah bro")
				childPage.ParentURL = append(childPage.ParentURL, page.GetURL())
			}
		}
	}
}

// WriteIndexed write the page data into an external file given a page map
func (page *Page) WriteIndexed(pages *map[string]*Page) {
	basePage := (*pages)[page.GetURL()]
	f, err := os.Create("../assets/spider_result.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	head := ""
	// head := "TITLE: " + basePage.GetTitle() + "\n" + basePage.GetURL() + "\n" + "DATE: " + basePage.GetLastModified() + ", " + basePage.GetSize() + "\n" + strings.Join(basePage.GetKeywords(), " ") + "\n" + strings.Join(basePage.GetChildrenURL(), "\n") + "\n"
	basePage.MakeChildren(pages)

	for _, child := range *pages {
		children := "----------------------------------------------\n" + "TITLE: " + child.GetTitle() + "\n" + child.GetURL() + "\n" + "DATE: " + child.GetLastModified() + ", " + child.GetSize() + "\n" + strings.Join(child.GetKeywords(), " ") + "\n" + strings.Join(child.GetChildrenURL(), "\n") + "\n"
		head += children

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
