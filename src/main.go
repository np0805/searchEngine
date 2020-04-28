package main

import (
	"fmt"
	"time"

	"./crawler"
)

// Crawl crawl a given url as its base url,
// returning a mapping of url --> page struct
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
	basePage.MakeChildren(&pagesMap)

	return pagesMap
}

func main() {
	const baseURL = "https://www.cse.ust.hk/"
	fmt.Println(time.Now()) // buat ngecek dia brp lama runnya

	pagesMap := Crawl(baseURL) // get the mapping of url --> page struct

	fmt.Println(time.Now())

	// contoh cara ngambil page dari map
	for _, page := range pagesMap {
		fmt.Println(page.GetTitle())
		fmt.Println(page.GetKeywords())
		break
	}
	/*
		contoh cara lain buat ngambil page dari map
		another := pagesMap["http://epublish.ust.hk/cgi-bin/eng/story.php?id=96&catid=97&keycode=88b7aae0ae45ddb0e6e000ee2682721a&token=17b43a00aeb0f8f8f08df16ae664909f"]
		fmt.Println(another.GetTitle())
	*/
}
