package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SignUp(w http.ResponseWriter, r *http.Request) *AppError {
	w.Header().Set("Content-Type", "application/json")

	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return DecodingRequestBodyError(err)
	}
	log.Println("user : ", user)

	foundUser, err := FindUserByEmail(user.Email)
	if err != nil {
		return DatabaseError(err, "")
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

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return EncodingResponseError(err, user)
	}

	return nil
}

type SignInParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func SignIn(w http.ResponseWriter, r *http.Request) *AppError {
	w.Header().Set("Content-Type", "application/json")

	var signInParams SignInParams
	if err := json.NewDecoder(r.Body).Decode(&signInParams); err != nil {
		return DecodingRequestBodyError(err)
	}
	log.Println("Sign In Params : ", signInParams)

	user, err := FindUserByEmail(signInParams.Email)
	if err != nil {
		return BadRequestError(err, "impossible to retreive the user")
	}
	log.Println("Found User: ", user)

	if user == nil {
		return BadRequestError(err, "No user with that email")
	}
	log.Println("User Found")

	res, err := CheckPassword(user.Password, signInParams.Password)
	if err != nil {
		return InternalError(err, "CheckPassword")
	}
	log.Println("Checkpassword : ", res)

	if !res {
		return BadRequestError(err, "Password don't match")
	}

	token, err := GetJWT(*user)
	if err != nil {
		return InternalError(err, "GetJWT")
	}
	log.Println("Token : ", token)

	response := SignInResponse{
		token,
		*user,
	}
	log.Println("Response : ", response)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return EncodingResponseError(err, response)
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

	r.Handle("/auth/signin", AppHandler(SignIn)).Methods(http.MethodPost)
	r.Handle("/auth/signup", AppHandler(SignUp)).Methods(http.MethodPost)
	r.HandleFunc("/auth/reset_password", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Forgot password")
	}).Methods(http.MethodPost)

	return r
}
