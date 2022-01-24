package main

import (
	"log"
	"net/http"
)

func OptionsMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	}
}

func RecoveryMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func MW(h http.HandlerFunc) http.HandlerFunc {
	return RecoveryMiddleware(OptionsMiddleware(h))
}
