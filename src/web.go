package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"./database"
	"./retrieval"
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
func main() {
	fmt.Println("Now Listening on 8000")
	database.OpenAllDb()
	http.HandleFunc("/", index)
	http.HandleFunc("/index", index)
	http.HandleFunc("/result", processinput)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
