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

// MakeChildren given a page, create its children page from the page's childrenURL and map it to the given map
func (page *Page) MakeChildren(pages *map[string]*Page) {
	for _, url := range page.GetChildrenURL() {
		childPage, ok := (*pages)[url]
		if !ok {
			childPage := Page{url, "", "", "", make([]string, 0), make([]string, 0), make([]string, 0)}
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
			childPage.ParentURL = append(childPage.ParentURL, page.GetURL())
		}

	}
}

// WriteIndexed write the page data into an external file given a page map
func (page *Page) WriteIndexed(pages *map[string]*Page) {
	basePage := (*pages)[page.GetURL()]
	// basePage.ExtractTitle()
	// basePage.ExtractLastModified()
	// basePage.ExtractWords()
	// basePage.ExtractSize()
	// basePage.ExtractLinks()
	f, err := os.Create("spider_result.txt")
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

// func main() {
// 	const baseURL = "https://www.cse.ust.hk/"
// 	fmt.Println(time.Now())

// 	// pages := make([]*Page, 0)
// 	// basePage := Page{baseURL, "", "", "", make([]string, 0), nil, make([]string, 0)}
// 	// pages = append(pages, &basePage)
// 	// WriteIndexed(&pages)

// 	pagesMap := make(map[string]*Page)
// 	basePage := Page{baseURL, "", "", "", make([]string, 0), nil, make([]string, 0)}

// 	basePage.ExtractTitle()
// 	basePage.ExtractLastModified()
// 	basePage.ExtractWords()
// 	basePage.ExtractSize()
// 	basePage.ExtractLinks()

// 	pagesMap[baseURL] = &basePage

// 	basePage.WriteIndexed(&pagesMap)

// 	// // contoh cara ngeindex dari map
// 	// another := pagesMap["http://epublish.ust.hk/cgi-bin/eng/story.php?id=96&catid=97&keycode=88b7aae0ae45ddb0e6e000ee2682721a&token=17b43a00aeb0f8f8f08df16ae664909f"]
// 	// fmt.Println(another.GetParentURL())
// 	// fmt.Println(pagesMap["http://epublish.ust.hk/cgi-bin/eng/story.php?id=96&catid=97&keycode=88b7aae0ae45ddb0e6e000ee2682721a&token=17b43a00aeb0f8f8f08df16ae664909f"])
// 	// fmt.Println(len(pagesMap))

// 	fmt.Println(time.Now())

// }

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
