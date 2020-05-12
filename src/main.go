package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"./crawler"
	"./database"
	"./pagerank"
	"./retrieval"
	"./stopstem"
)

var tpl *template.Template

type ResultPage struct {
	Id           int64    `json:"Id"`
	Score        float64  `json:"Score"`
	Title        string   `json:"Title"`
	Url          string   `json:"Url"`
	LastModified string   `json:"LastModified"`
	PageSize     string   `json:"PageSize"`
	Keywords     []string `json:"Keywords"`
	Parents      []string `json:"Parents"`
	Children     []string `json:"Children"`
}

type TheResult struct {
	ResultPages []ResultPage `json:"result"`
}

// func init() {
//   tpl = template.Must(template.ParseGlob("templates/*html"))
// }

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("../index.html")
		t.Execute(w, nil)
	} else {

		//Parsing contents of the form
		r.ParseForm()

		start := time.Now()
		//Submitting query to the search engine
		elapsed := time.Since(start)
		fmt.Println("Query took ", elapsed)

		http.Redirect(w, r, "/result", http.StatusSeeOther)
	}
	// tpl.ExecuteTemplate(w, "index.html", nil)
}
func processinput(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	//   http.Redirect(w, r, "/", http.StatusSeeOther)
	//   return
	// }

	//create json files of the result
	fmt.Println(time.Now())
	retrieval.RetrievalFunction(r.FormValue("searchInput"))
	fmt.Println(time.Now())

	//extract the json files
	jsonFile, err := os.Open("search_output.json")

	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println("connected to search_output.json")
	// fmt.Print(jsonFile)
	defer jsonFile.Close()
	// fmt.Print(jsonFile)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var results TheResult
	// var results map[string]interface{}
	fmt.Println("STEP 2")
	json.Unmarshal([]byte(byteValue), &results)
	// fmt.Print(results)
	fmt.Println("STEP 4")
	var resultString string
	for i, rel := range results.ResultPages {
		if i == 50 {
			break
		}
		resultString = resultString + "<p style=\"font-size: 10px\">Score: " + database.FloatToString(rel.Score) + "</p>"
		resultString = resultString + "<strong><a href=\"" + rel.Url + "\">" + rel.Title + "</a></strong><br>"
		resultString = resultString + "URL: " + rel.Url + "<br>"
		resultString = resultString + rel.LastModified + " | " + rel.PageSize + "<br><br>"
		resultString = resultString + "<strong>5 top keywords:</strong> <br>" + database.SliceToString(rel.Keywords) + "<br>"
		resultString = resultString + "<Strong>Parents: </Strong><br>"
		for j, par := range rel.Parents {
			if j == 5 {
				break
			}
			resultString = resultString + par + "<br>"
		}
		resultString = resultString + "<strong>Children: </strong><br>"
		for j, chil := range rel.Children {
			if j == 5 {
				break
			}
			resultString = resultString + chil + "<br>"
		}
		resultString = resultString + "<br>"
	}
	// fmt.Println(resultString)

	// // tpl.ExecuteTemplate(w, "result.html", results)

	// tpl.Execute(w, map[string]interface{}{
	//    ".": html,
	// })
	t, _ := template.ParseFiles("../result.html")
	html := template.HTML(resultString)
	t.Execute(w, map[string]interface{}{
		"Body": html,
	})
}

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
	basePage.WriteIndexed(&pagesMap)

	return pagesMap
}

func main() {
	const baseURL = "https://www.cse.ust.hk/"
	fmt.Println(time.Now()) // buat ngecek dia brp lama runnya

	pagesMap := Crawl(baseURL) // get the mapping of url --> page struct
	fmt.Println("Len of map %v", len(pagesMap))
	fmt.Println(time.Now()) // buat ngecek dia brp lama runnya

	pagerank.CalculatePageRank(0.85, &pagesMap)

	// // contoh cara ngambil page dari map
	// for _, page := range pagesMap {
	// 	fmt.Println(page.GetURL(), page.GetPageRank())
	// }

	// mapAwal := pagesMap["https://www.cse.ust.hk/admin/people/staff/"]
	// fmt.Println(mapAwal.GetTitle())
	// fmt.Println(mapAwal.GetKeywords())
	/*
	  contoh cara lain buat ngambil page dari map
	  another := pagesMap["http://epublish.ust.hk/cgi-bin/eng/story.php?id=96&catid=97&keycode=88b7aae0ae45ddb0e6e000ee2682721a&token=17b43a00aeb0f8f8f08df16ae664909f"]
	  fmt.Println(another.GetTitle())
	*/
	fmt.Println("-------------------------------------------------------")
	stopstem.InputStopWords()
	newMap := stopstem.StemThemAll(&pagesMap)
	fmt.Println(len(newMap))
	fmt.Println(time.Now()) // buat ngecek dia brp lama runnya
	// for _, page := range newMap {
	// 	fmt.Println("pageGetURL: ", page.GetURL())
	// 	fmt.Println("getPageRank: ", page.GetPageRank())
	// }
	database.OpenAllDb()
	database.ParseAllPages(&newMap)
	//database.PrintPageIdDb()
	fmt.Println(time.Now())
	// mapAkhir := newMap["https://www.cse.ust.hk/admin/people/staff/"]
	// fmt.Println(mapAkhir.GetTitle())
	// fmt.Println(mapAkhir.GetKeywords())
}
