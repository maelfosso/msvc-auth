package controllers

import "guitou.cm/msvc/auth/models"

type SignInParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type ForgetPasswordParams struct {
	Email string `json:"email"`
}
