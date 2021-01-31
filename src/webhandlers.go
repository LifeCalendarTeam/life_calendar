package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

const MaxProportion = 100

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

		// Making `passwordHash` a slice of strings instead of just a string because we can only scan from SQL to a
		// slice. In fact the length of `passwordHash` will always be either 0 (if user doesn't exist) or 1
		passwordHash := make([]string, 0, 1)

		panicIfError(db.Select(&passwordHash, "SELECT password_hash FROM users WHERE id=$1", person.UserID))
		if len(passwordHash) == 0 {
			http.Error(w, "There is no user with the given id", http.StatusForbidden)
			return
		}
		bcryptErr := bcrypt.CompareHashAndPassword([]byte(passwordHash[0]), []byte(person.Password))
		if bcryptErr == bcrypt.ErrMismatchedHashAndPassword {
			http.Error(w, "Looks like your id or password is incorrect", http.StatusForbidden)
		} else if bcryptErr == nil {
			session, _ := cookieStorage.Get(r, "session")
			session.Values["id"] = person.UserID
			session.Values["expires"] = time.Now().Add(24 * time.Hour).Unix()
			panicIfError(session.Save(r, w))

			http.Redirect(w, r, "..", 302)
		} else {
			panic(bcryptErr)
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
	for proportionColorIdx := range proportionsAndColors {
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
		for idx := range ProportionedTotalColor {
			absoluteColor := ProportionedTotalColor[idx] * MaxProportion / TotalProportion
			ans[idx] = int(math.Round(absoluteColor))
		}
	}

	return ans, nil
}

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
		writeJSON(w, js, http.StatusBadRequest)
		return
	}

	// TODO: the two lines below are formatted automatically this way. I wonder if it is possible to write it better
	if len(r.Form["activity_type"]) != len(r.Form["activity_proportion"]) ||
		len(r.Form["emotion_type"]) != len(r.Form["emotion_proportion"]) {

		js, err := json.Marshal(map[string]interface{}{"ok": false,
			"error": "Lengths of types and proportions of both activities and emotions must be equal " +
				"correspondingly",
			"error_type": "types_and_proportions_lengths"})
		panicIfError(err)
		writeJSON(w, js, http.StatusBadRequest)
		return
	}

	activitiesEmotionsTypes := append(r.Form["activity_type"], r.Form["emotion_type"]...)
	activitiesEmotionsProportions := append(r.Form["activity_proportion"], r.Form["emotion_proportion"]...)

	tx, err := db.Begin()
	panicIfError(err)
	defer func() {
		_ = tx.Rollback()
	}()

	var dayID int
	err = tx.QueryRow("INSERT INTO days(user_id, date) VALUES ($1, $2) RETURNING id", userID, date).Scan(&dayID)
	if pgErr, ok := err.(*pq.Error); ok {
		if pgErr.Constraint == "days_user_id_date_key" {
			js, err := json.Marshal(map[string]interface{}{"ok": false, "error": "The user has a day with the date",
				"error_type": "day_already_exists"})
			panicIfError(err)
			writeJSON(w, js, http.StatusPreconditionFailed)
			return
		}
	}
	panicIfError(err)
	for idx := range activitiesEmotionsTypes {
		proportion, err := strconv.Atoi(activitiesEmotionsProportions[idx])
		if err != nil {
			js, err := json.Marshal(map[string]interface{}{"ok": false, "error": "Activity/emotion proportion must " +
				"be an integer, not \"" + activitiesEmotionsProportions[idx] + "\"",
				"error_type": "incorrect_proportion"})
			panicIfError(err)
			writeJSON(w, js, http.StatusBadRequest)
			return
		}
		if proportion < 1 || proportion > MaxProportion {
			js, err := json.Marshal(map[string]interface{}{"ok": false, "error": "The proportion must be from 1 to " +
				strconv.Itoa(MaxProportion), "error_type": "incorrect_proportion"})
			panicIfError(err)
			writeJSON(w, js, http.StatusPreconditionFailed)
			return
		}

		_, err = tx.Exec("INSERT INTO activities_and_emotions(type_id, day_id, proportion) VALUES ($1, $2, $3)",
			activitiesEmotionsTypes[idx], dayID, activitiesEmotionsProportions[idx])
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Constraint == "activities_and_emotions_type_id_fkey" ||
				pgErr.Constraint == "type_owned_by_correct_user_check" {

				js, err := json.Marshal(map[string]interface{}{"ok": false, "error": "The user doesn't have the " +
					"activity/emotion with type " + activitiesEmotionsTypes[idx], "error_type": "incorrect_type"})
				panicIfError(err)
				writeJSON(w, js, http.StatusPreconditionFailed)
				return
			}

			// activity/emotion's `type_id` is not an `int`:
			// (note that this error isn't about `proportion`, because it was earlier successfully converted to `int`)
			if pgErr.Code.Name() == "invalid_text_representation" {
				js, err := json.Marshal(map[string]interface{}{"ok": false, "error": "Activity/emotion type id must " +
					"be an integer, not \"" + activitiesEmotionsTypes[idx] + "\"", "error_type": "incorrect_type"})
				panicIfError(err)
				writeJSON(w, js, http.StatusBadRequest)
				return
			}

			// TODO: also implement and document the handling of the case when the request contained activity/emotion
			// multiple times
		}
		panicIfError(err)
	}
	err = tx.Commit()
	panicIfError(err)

	js, err := json.Marshal(map[string]interface{}{"ok": true, "id": dayID})
	panicIfError(err)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	panicIfError(err)
}

func HandleAPIDaysBrief(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStorage.Get(r, "session")
	if dropAPIRequestIfUnauthorized(session, w) {
		return
	}

	days := make([]briefDay, 0)
	panicIfError(db.Select(&days, "SELECT id, date FROM days WHERE user_id=$1", session.Values["id"]))

	// Retrieving average color:
	for dayIdx := range days {
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
	writeJSON(w, js, http.StatusOK)
}

func HandleAPIDaysID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	// `err` should never occur because Gorilla should have rejected the request before calling the handler if `id` is
	// not an int
	panicIfError(err)

	session, _ := cookieStorage.Get(r, "session")
	if dropAPIRequestIfUnauthorized(session, w) {
		return
	}

	b := make([]bool, 0)
	panicIfError(db.Select(&b, "SELECT EXISTS(SELECT 1 FROM days WHERE id=$1 AND user_id=$2)", id,
		session.Values["id"]))
	dayExists := b[0]

	var js []byte
	var status = http.StatusOK

	if dayExists {
		if r.Method == "DELETE" {
			_, err = db.Exec("DELETE FROM days WHERE id=$1", id)
			if err == nil {
				js, err = json.Marshal(map[string]interface{}{"ok": true})
			}
		} else {
			activitiesAndEmotions := make([]activityOrEmotionWithType, 0)
			panicIfError(db.Select(&activitiesAndEmotions, "SELECT type_id, proportion, "+
				"(SELECT activity_or_emotion FROM types_of_activities_and_emotions WHERE id=type_id) FROM "+
				"activities_and_emotions WHERE day_id=$1", id))

			activities := make([]ActivityOrEmotion, 0, len(activitiesAndEmotions))
			emotions := make([]ActivityOrEmotion, 0, len(activitiesAndEmotions))
			for _, entity := range activitiesAndEmotions {
				entity.DayID = 0 // Hide the field in JSON
				if entity.EntityType == EntityTypeActivity {
					activities = append(activities, entity.ActivityOrEmotion)
				} else {
					emotions = append(emotions, entity.ActivityOrEmotion)
				}
			}

			js, err = json.Marshal(map[string]interface{}{"ok": true, "activities": activities, "emotions": emotions})
		}
	} else {
		js, err = json.Marshal(map[string]interface{}{"ok": false, "error": "Day does not exist"})
		status = http.StatusNotFound
	}

	panicIfError(err)
	writeJSON(w, js, status)
}

func main() {
	ui := mux.NewRouter()
	ui.Path("/").Methods("GET").HandlerFunc(HandleRoot)
	ui.Path("/login").Methods("GET", "POST").HandlerFunc(HandleLogin)
	ui.Path("/logout").Methods("GET").HandlerFunc(HandleLogout)

	api := mux.NewRouter()
	api.Path("/api/days").Methods("POST").HandlerFunc(HandleAPIDays)
	api.Path("/api/days/brief").Methods("GET").HandlerFunc(HandleAPIDaysBrief)
	api.Path("/api/days/{id:[0-9]+}").Methods("GET", "DELETE").HandlerFunc(HandleAPIDaysID)

	final := http.NewServeMux()
	final.Handle("/", UIPanicHandlerMiddleware(ui))
	final.Handle("/api/", APIPanicHandlerMiddleware(api))
	// TODO: Authorization check middleware

	listenAddr := "localhost:4000"
	fmt.Println("Listening at http://" + listenAddr)
	panic(http.ListenAndServe(listenAddr, final))
}
