package main

import (
	"fmt"
	"time"

	"./crawler"
)

// Crawl crawl a given url as its base url, returning a mapping of url --> page struct
func Crawl(baseURL string) map[string]*crawler.Page {
	pagesMap := make(map[string]*crawler.Page)
	basePage := crawler.Page{
		URL:          baseURL,
		Title:        "",
		LastModified: "",
		PageSize:     "",
		Keywords:     make([]string, 0),
		ParentURL:    nil,
		ChildrenURL:  make([]string, 0)}

	basePage.ExtractTitle()
	basePage.ExtractLastModified()
	basePage.ExtractWords()
	basePage.ExtractSize()
	basePage.ExtractLinks()
	pagesMap[baseURL] = &basePage

	return pagesMap
}

func main() {
	const baseURL = "https://www.cse.ust.hk/"
	fmt.Println(time.Now())

	pagesMap := Crawl(baseURL)
	basePage := pagesMap[baseURL]

	basePage.WriteIndexed(&pagesMap)

	// // contoh cara ngeindex dari map
	// another := pagesMap["http://epublish.ust.hk/cgi-bin/eng/story.php?id=96&catid=97&keycode=88b7aae0ae45ddb0e6e000ee2682721a&token=17b43a00aeb0f8f8f08df16ae664909f"]
	// fmt.Println(another.GetParentURL())
	// fmt.Println(pagesMap["http://epublish.ust.hk/cgi-bin/eng/story.php?id=96&catid=97&keycode=88b7aae0ae45ddb0e6e000ee2682721a&token=17b43a00aeb0f8f8f08df16ae664909f"])
	// fmt.Println(len(pagesMap))

	fmt.Println(time.Now())

}
