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

// writeJSON marshals the given data (`data`) with json, writes the result to the `http.ResponseWriter` (`w`) ans uses
// the `status` HTTP status code. Panics if either of `json.Marshal` or `w.Write` has failed
func writeJSON(w http.ResponseWriter, data map[string]interface{}, status int) {
	js, err := json.Marshal(data)
	panicIfError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	panicIfError(err)
}

// dropAPIRequestIfUnauthorized checks if the session object is new (`session.IsNew`) and sends an error json as a
// response if it is. Does nothing otherwise. Return-value is a `bool`, which is `true` if the user is unauthorized
// (and, therefore, the error response was written) `false` if the session is not new
func dropAPIRequestIfUnauthorized(session *sessions.Session, w http.ResponseWriter) bool {
	if session.IsNew {
		writeJSON(w,
			map[string]interface{}{"ok": false, "error": "You must be authorized to call this method"},
			http.StatusUnauthorized)
		return true
	}
	return false
}
