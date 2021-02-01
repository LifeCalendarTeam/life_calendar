package main

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

// panicIfError calls `log.Panic(err)` if `err` is not `nil`
func panicIfError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// writeJSON writes the given marshalled json (`js`) to an `http.ResponseWriter` (`w`) and uses the `status` HTTP status
// code. Panics if `w.Write` had failed
func writeJSON(w http.ResponseWriter, js []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(js)
	panicIfError(err)
}

// dropAPIRequestIfUnauthorized checks if the session object is new (`session.IsNew`) and sends an error json as a
// response if it is. Does nothing otherwise. Return-value is a `bool`, which is `true` if the user is unauthorized
// (and, therefore, the error response was written) `false` if the session is not new
func dropAPIRequestIfUnauthorized(session *sessions.Session, w http.ResponseWriter) bool {
	if session.IsNew {
		js, err := json.Marshal(map[string]interface{}{"ok": false,
			"error": "You must be authorized to call this method"})
		panicIfError(err)
		writeJSON(w, js, http.StatusUnauthorized)
		return true
	}
	return false
}
