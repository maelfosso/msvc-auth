package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"guitou.cm/msvc/auth/db"
	"guitou.cm/msvc/auth/models"
	"guitou.cm/msvc/auth/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) *AppError {
	w.Header().Set("Content-Type", "application/json")

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return DecodingRequestBodyError(err)
	}
	log.Println("user : ", user)

	foundUser, err := db.FindUserByEmail(user.Email)
	if err != nil {
		return DatabaseError(err, "")
	}
	log.Println("Found User: ", foundUser)

	if foundUser != nil {
		return BadRequestError(err, "User already exists")
	}
	log.Println("No user found")

	id, err := db.SaveUser(user)
	if err != nil {
		return DatabaseError(err, "saving user failed")
	}
	user.ID = id
	log.Println("User saved: ", user)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return EncodingResponseError(err, user)
	}

	return nil
}

func SignIn(w http.ResponseWriter, r *http.Request) *AppError {
	w.Header().Set("Content-Type", "application/json")

	var signInParams SignInParams
	if err := json.NewDecoder(r.Body).Decode(&signInParams); err != nil {
		return DecodingRequestBodyError(err)
	}
	log.Println("Sign In Params : ", signInParams)

	user, err := db.FindUserByEmail(signInParams.Email)
	if err != nil {
		return DatabaseError(err, "impossible to retreive the user")
	}
	log.Println("Found User: ", user)

	if user == nil {
		return BadRequestError(err, "No user with that email")
	}
	log.Println("User Found")

	res := utils.CheckPassword(user.Password, signInParams.Password)
	if !res {
		return BadRequestError(err, "Password don't match")
	}

	token, err := utils.GetJWT(*user)
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

func Check(w http.ResponseWriter, r *http.Request) *AppError {

	return nil
}
