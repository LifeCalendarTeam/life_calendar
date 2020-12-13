package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// APIPanicHandlerMiddleware is an HTTP middleware, handling panic and sending the error json as a response in case of
// errors
func APIPanicHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)

				js, err := json.Marshal(map[string]interface{}{"ok": false, "error": fmt.Sprint(err)})
				if err != nil {
					log.Println(err)
					http.Error(w, "{\"ok\":false,\"error\":\"unknown\"}", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				http.Error(w, string(js), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// UIPanicHandlerMiddleware is an HTTP middleware, handling panic and sending the user-readable error page as a response
// in case of errors
func UIPanicHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)

				w.WriteHeader(http.StatusInternalServerError)
				if err := tmpl.ExecuteTemplate(w, "500", fmt.Sprint(err)); err != nil {
					log.Println(err)

					// I don't think we need to handle the following error, because what can we do if it occurs?
					_, _ = w.Write([]byte("It failed :("))
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
