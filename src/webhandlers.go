package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var tmpl *template.Template
var cookieStorage *sessions.CookieStore
var requestsDecoder = schema.NewDecoder()
var MaxProportion = 100.

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

	var cookieKey []byte
	if cookieKeyStr, exists := os.LookupEnv("cookie_key"); exists {
		cookieKey = []byte(cookieKeyStr)
	} else {
		cookieKey = securecookie.GenerateRandomKey(32)
	}
	cookieStorage = sessions.NewCookieStore(cookieKey)
}

type loginForm struct {
	UserId   int    `schema:"user_id,required"`
	Password string `schema:"password,required"`
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStorage.Get(r, "session")

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
			session, _ := cookieStorage.Get(r, "session")
			session.Values["id"] = person.UserId
			session.Values["expires"] = time.Now().Add(24 * time.Hour).Unix()
			panicIfError(session.Save(r, w))

			http.Redirect(w, r, "..", 302)
		} else {
			http.Error(w, "Looks like your login or password is incorrect", http.StatusInternalServerError)
		}
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStorage.Get(r, "session")
	session.Options.MaxAge = -1
	panicIfError(session.Save(r, w))

	http.Redirect(w, r, "..", 302)
}

func handleApiDaysBrief(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStorage.Get(r, "session")
	if session.IsNew {
		http.Error(w, "You should be authorized to call this method", http.StatusUnauthorized)
		return
	}

	type BriefDay struct {
		DayId        int       `db:"id" json:"id"`
		Date         time.Time `db:"date" json:"date"`
		AverageColor [3]int    `json:"average_color"`
	}
	days := make([]BriefDay, 0)
	panicIfError(db.Select(&days, "SELECT id, date FROM days WHERE user_id=$1", session.Values["id"]))

	// Retrieving average color:
	for dayIdx, _ := range days {
		type ProportionAndColor struct {
			Proportion float64 `db:"proportion"`
			Color      string  `db:"color"`
		}
		day := &days[dayIdx]

		colorsProportions := make([]ProportionAndColor, 0)
		panicIfError(db.Select(&colorsProportions,
			"SELECT CAST(proportion AS FLOAT), (SELECT color FROM types_of_emotions WHERE id = type_id) FROM emotions",
		))
		panicIfError(db.Select(&colorsProportions,
			"SELECT CAST(proportion AS FLOAT), (SELECT color FROM types_of_activities WHERE id = type_id) "+
				"FROM activities",
		))

		TotalProportion := 0.
		ProportionedTotalColor := [...]float64{0, 0, 0}
		for proportionColorIdx, _ := range colorsProportions {
			proportionColor := &colorsProportions[proportionColorIdx]

			TotalProportion += proportionColor.Proportion
			colorAsStringArray := strings.Split(proportionColor.Color, ",")

			if len(colorAsStringArray) != 3 {
				panic(errors.New("exactly three colors must be in the `color` field of the database"))
			}

			for idx, val := range colorAsStringArray {
				col, err := strconv.Atoi(val)
				panicIfError(err)
				ProportionedTotalColor[idx] += float64(col) * proportionColor.Proportion / MaxProportion
			}
		}

		for idx, _ := range ProportionedTotalColor {
			absoluteColor := ProportionedTotalColor[idx] * MaxProportion / TotalProportion
			day.AverageColor[idx] = int(math.Round(absoluteColor))
		}
	}

	js, err := json.Marshal(days)
	panicIfError(err)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot).Methods("GET")
	r.HandleFunc("/login", handleLogin).Methods("GET", "POST")
	r.HandleFunc("/logout", handleLogout).Methods("GET")

	r.HandleFunc("/api/days/brief", handleApiDaysBrief).Methods("GET")
	//r.HandleFunc("/api/days/{id:[0-9]+}")

	listenAddr := "localhost:4000"
	fmt.Println("Listening at http://" + listenAddr)
	panic(http.ListenAndServe(listenAddr, r))
}
