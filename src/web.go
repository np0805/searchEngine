package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"./database"
	"./retrieval"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*html"))
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}
func processinput(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	database.OpenAllDb()
	fmt.Println(time.Now())
	t := r.FormValue("searchInput")
	retrieval.RetrievalFunction(t)
	//typeresult := reflect.TypeOf(resultmap)
	//teststring := resultmap[0].GetTitle()
	//fmt.Println(teststring)
	d := struct {
		Time   time.Time
		String string
	}{
		Time:   time.Now(),
		String: t,
	}

	tpl.ExecuteTemplate(w, "result.html", d)
	database.CloseAlldb()
}
func main() {
	fmt.Println("Now Listening on 8000")
	http.HandleFunc("/", index)
	http.HandleFunc("/result", processinput)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
