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

func loadHTMLTemplates() (*template.Template, error) {
	tmpl := template.New("HTML templates")
	return tmpl.ParseFiles("src/templates/index.html", "src/templates/login.html", "src/templates/500.html")
}

func init() {
	var err error

	tmpl, err = loadHTMLTemplates()
	panicIfError(err)

	var cookieKey []byte
	if cookieKeyStr, exists := os.LookupEnv("cookie_key"); exists {
		cookieKey = []byte(cookieKeyStr)
	} else {
		cookieKey = securecookie.GenerateRandomKey(32)
	}
	cookieStorage = sessions.NewCookieStore(cookieKey)
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStorage.Get(r, "session")

	panicIfError(tmpl.ExecuteTemplate(w, "/", !session.IsNew))
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		panicIfError(tmpl.ExecuteTemplate(w, "GET /login", nil))
	} else {
		panicIfError(r.ParseForm())

		var person loginForm

		if err := requestsDecoder.Decode(&person, r.Form); err != nil {
			http.Error(w, "Looks like you have sent incorrect data", http.StatusBadRequest)
		}

		if true { // Password must be checked here!!
			session, _ := cookieStorage.Get(r, "session")
			session.Values["id"] = person.UserID
			session.Values["expires"] = time.Now().Add(24 * time.Hour).Unix()
			panicIfError(session.Save(r, w))

			http.Redirect(w, r, "..", 302)
		} else {
			http.Error(w, "Looks like your login or password is incorrect", http.StatusForbidden)
		}
	}
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStorage.Get(r, "session")
	session.Options.MaxAge = -1
	panicIfError(session.Save(r, w))

	http.Redirect(w, r, "..", 302)
}

func getAverageColor(proportionsAndColors []proportionAndColor) ([3]int, error) {
	TotalProportion := 0.
	ProportionedTotalColor := [...]float64{0, 0, 0}
	for proportionColorIdx, _ := range proportionsAndColors {
		proportionColor := &proportionsAndColors[proportionColorIdx]

		TotalProportion += proportionColor.Proportion
		colorAsStringArray := strings.Split(proportionColor.Color, ",")

		if len(colorAsStringArray) != 3 {
			return [3]int{0, 0, 0}, errors.New("exactly three colors must be in the `color` field of the database")
		}

		for idx, val := range colorAsStringArray {
			col, err := strconv.Atoi(val)
			panicIfError(err)
			ProportionedTotalColor[idx] += float64(col) * proportionColor.Proportion / MaxProportion
		}
	}

	ans := [3]int{0, 0, 0}

	if TotalProportion != 0 {
		for idx, _ := range ProportionedTotalColor {
			absoluteColor := ProportionedTotalColor[idx] * MaxProportion / TotalProportion
			ans[idx] = int(math.Round(absoluteColor))
		}
	}

	return ans, nil
}

// TODO: not forget to replace all the json fuck with `writeJSON` when fixing merge conflicts
func HandleAPIDays(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStorage.Get(r, "session")
	// Panics if user is not authorized. Will be fixed with the appropriate middleware
	userID := session.Values["id"].(int)

	panicIfError(r.ParseForm())

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		js, err := json.Marshal(map[string]interface{}{"ok": false, "error": "Unable to parse date",
			"error_type": "incorrect_date"})
		panicIfError(err)
		http.Error(w, string(js), http.StatusPreconditionFailed)
		return
	}

	// TODO: the two lines below are formatted automatically this way. I wonder if it is possible to write it better
	if len(r.Form["activity_type"]) != len(r.Form["activity_proportion"]) ||
		len(r.Form["emotion_type"]) != len(r.Form["emotion_proportion"]) {

		js, err := json.Marshal(map[string]interface{}{"ok": false,
			"error":      "Lengths of types and proportions of both activities and emotions must be equal correspondingly",
			"error_type": "types_and_proportions_lengths"})
		panicIfError(err)
		http.Error(w, string(js), http.StatusBadRequest)
		return
	}

	activitiesEmotionsTypes := append(r.Form["activity_type"], r.Form["emotion_type"]...)
	activitiesEmotionsProportions := append(r.Form["activity_proportion"], r.Form["emotion_proportion"]...)

	// TODO: all the following must be under one transaction!
	res, err := db.Exec("INSERT INTO days(user_id, date) VALUES ($1, $2)", userID, date)
	panicIfError(err) // TODO: probably this is a requester's mistake! Should return an appropriate error then
	dayID, err := res.LastInsertId()
	panicIfError(err) // TODO: figure out if this can ever happen with PostgreSQL. Probably we can omit the check
	for idx := range activitiesEmotionsTypes {
		// TODO: check if `type_id` belongs to the correct user. If not, should return 412. Btw, I believe this should
		// be a PostgreSQL constraint
		_, err = db.Exec("INSERT INTO activities_and_emotions(type_id, day_id, proportion) VALUES ($1, $2, $3)",
			activitiesEmotionsTypes[idx], dayID, activitiesEmotionsProportions[idx])
	}

	js, err := json.Marshal(map[string]interface{}{"ok": true, "id": dayID})
	panicIfError(err)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	panicIfError(err)
}

func HandleAPIDaysBrief(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStorage.Get(r, "session")
	if session.IsNew {
		http.Error(w, "You should be authorized to call this method", http.StatusUnauthorized)
		return
	}

	days := make([]briefDay, 0)
	panicIfError(db.Select(&days, "SELECT id, date FROM days WHERE user_id=$1", session.Values["id"]))

	// Retrieving average color:
	for dayIdx, _ := range days {
		day := &days[dayIdx]

		colorsProportions := make([]proportionAndColor, 0)
		panicIfError(db.Select(&colorsProportions,
			"SELECT CAST(proportion AS FLOAT), (SELECT color FROM types_of_activities_and_emotions "+
				"WHERE id=type_id) FROM activities_and_emotions WHERE day_id=$1", day.DayID,
		))

		var err error
		day.AverageColor, err = getAverageColor(colorsProportions)
		panicIfError(err)
	}

	js, err := json.Marshal(map[string]interface{}{"ok": true, "days": days})
	panicIfError(err)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	panicIfError(err)
}

func main() {
	ui := mux.NewRouter()

	ui.HandleFunc("/", HandleRoot).Methods("GET")
	ui.HandleFunc("/login", HandleLogin).Methods("GET", "POST")
	ui.HandleFunc("/logout", HandleLogout).Methods("GET")

	api := mux.NewRouter()

	api.HandleFunc("/api/2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi..."))
	})
	api.HandleFunc("/api/days", HandleAPIDays).Methods("POST")
	api.HandleFunc("/api/days/brief", HandleAPIDaysBrief).Methods("GET")
	//api.HandleFunc("/api/days/{id:[0-9]+}")

	final := http.NewServeMux()
	final.Handle("/", UIPanicHandlerMiddleware(ui))
	final.Handle("/api/", APIPanicHandlerMiddleware(api))
	// TODO: Authorization check middleware

	listenAddr := "localhost:4000"
	fmt.Println("Listening at http://" + listenAddr)
	panic(http.ListenAndServe(listenAddr, final))
}
