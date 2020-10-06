package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
)

var tmpl *template.Template

func loadHtmlTemplates() (*template.Template, error) {
	if err := os.Chdir("src/templates"); err != nil {
		return nil, err
	}
	defer os.Chdir("../..")

	tmpl := template.New("HTML templates")
	return tmpl.ParseFiles("index.html")
}

func init() {
	var err error

	tmpl, err = loadHtmlTemplates()
	panicIfError(err)
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	panicIfError(tmpl.ExecuteTemplate(w, "/", nil))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot).Methods("GET")

	listenAddr := "localhost:4000"
	fmt.Println("Listening at http://" + listenAddr)
	panic(http.ListenAndServe(listenAddr, r))
}
