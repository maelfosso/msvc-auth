package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"guitou.cm/msvc/auth/controllers"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("I'm alive")
	})

	r.Handle("/auth/signin", AppHandler(controllers.SignIn)).Methods(http.MethodPost)
	r.Handle("/auth/signup", AppHandler(controllers.SignUp)).Methods(http.MethodPost)
	r.Handle("/auth/check", AppHandler(controllers.Check)).Methods(http.MethodPost)

	r.HandleFunc("/auth/password/forgot", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Forgot password")
	}).Methods(http.MethodPost)
	r.HandleFunc("/auth/password/reset", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Forgot password")
	}).Methods(http.MethodPost)

	return r
}
