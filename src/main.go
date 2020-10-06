package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"os"
	"time"
)

var tmpl *template.Template
var requestsDecoder = schema.NewDecoder()
var storage = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

func loadHtmlTemplates() (*template.Template, error) {
	if err := os.Chdir("src/templates"); err != nil {
		return nil, err
	}
	defer os.Chdir("../..")

	tmpl := template.New("HTML templates")
	return tmpl.ParseFiles("index.html", "login.html")
}

func init() {
	var err error

	tmpl, err = loadHtmlTemplates()
	panicIfError(err)
}

type loginForm struct {
	UserId   int    `schema:"user_id,required"`
	Password string `schema:"password,required"`
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	session, _ := storage.Get(r, "session")

	panicIfError(tmpl.ExecuteTemplate(w, "/", !session.IsNew))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		panicIfError(tmpl.ExecuteTemplate(w, "GET /login", nil))
	} else {
		panicIfError(r.ParseForm())

		var person loginForm

		if err := requestsDecoder.Decode(&person, r.Form); err != nil {
			http.Error(w, "Looks like you have sent incorrect data", http.StatusInternalServerError)
		}

		if true { // Password must be checked here!!
			session, _ := storage.Get(r, "session")
			session.Values["id"] = person.UserId
			session.Values["expires"] = time.Now().Add(24 * time.Hour).Unix()
			panicIfError(session.Save(r, w))

			http.Redirect(w, r, "..", 302)
		} else {
			http.Error(w, "Looks like your login or password is incorrect", http.StatusInternalServerError)
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot).Methods("GET")
	r.HandleFunc("/login", handleLogin).Methods("GET", "POST")

	listenAddr := "localhost:4000"
	fmt.Println("Listening at http://" + listenAddr)
	panic(http.ListenAndServe(listenAddr, r))
}
