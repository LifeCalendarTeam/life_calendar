package main

import (
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
