package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"guitou.cm/msvc/auth/controllers"
)

type AppHandler func(http.ResponseWriter, *http.Request) *controllers.AppError

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Message, err.StatusCode)
	}
}

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

	r.Handle("/auth/password/forgot", AppHandler(controllers.ForgetPassword)).Methods(http.MethodPost)
	r.Handle("/auth/password/forgot/verify", AppHandler(controllers.VerifyResetPasswordToken)).Methods(http.MethodPost)
	r.Handle("/auth/password/reset", AppHandler(controllers.ResettingPassword)).Methods(http.MethodPost)

	return r
}
