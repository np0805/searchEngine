package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"../src/retrieval"
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
	resultmap := retrieval.RetrievalFunction(r.FormValue("searchInput"))
	//typeresult := reflect.TypeOf(resultmap)
	/* d := struct {
		String reflect.Type
	}{
		String: typeresult}
	*/
	tpl.ExecuteTemplate(w, "result.html", resultmap)
}
func main() {
	fmt.Println("Now Listening on 8000")
	http.HandleFunc("/", index)
	http.HandleFunc("/result", processinput)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
