package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
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
	uInput := r.FormValue("searchInput")

	d := struct {
		String string
	}{
		String: uInput}

	tpl.ExecuteTemplate(w, "result.html", d)
}
func main() {
	fmt.Println("Now Listening on 8000")
	http.HandleFunc("/", index)
	http.HandleFunc("/result", processinput)
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
