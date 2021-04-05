package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SignIn(w http.ResponseWriter, r *http.Request) {

}

func SignUp(w http.ResponseWriter, r *http.Request) *AppError {
	w.Header().Set("Content-Type", "application/json")

	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return DecodeBodyDataError(err)
	}
	log.Println("user : ", user)

	foundUser, err := FindUserByEmail(user.Email)
	if err != nil {
		return BadRequestError(err, "")
	}
	log.Println("Found User: ", foundUser)

	if foundUser != nil {
		return BadRequestError(err, "User already exists")
	}
	log.Println("No user found")

	if err = user.Save(); err != nil {
		return DatabaseError(err, "saving user failed")
	}
	log.Println("User saved: ", user)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return EncodingResponseError(err, user)
	}

	return nil
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {

}

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("I'm alive")
	})
	// r.Handle()

	r.HandleFunc("/auth/signin", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Sign In")
	}).Methods(http.MethodPost)
	r.Handle("/auth/signup", AppHandler(SignUp)).Methods(http.MethodPost)
	// http.Handler
	r.HandleFunc("/auth/reset_password", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Forgot password")
	}).Methods(http.MethodPost)

	return r
}
