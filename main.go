package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
	"guitou.cm/msvc/auth/db"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		log.Printf(
			"%s\t%s\t%s\t",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
		next.ServeHTTP(w, r)
	})
}

func main() {
	db.OpenDB()

	r := NewRouter()
	r.Use(loggingMiddleware)

	handler := cors.AllowAll().Handler(r)

	log.Println("Listen on port :6000")
	log.Fatal(http.ListenAndServe(":6000", handler))
}
