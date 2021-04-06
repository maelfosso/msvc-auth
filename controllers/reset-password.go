package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"guitou.cm/msvc/auth/db"
)

func ForgetPassword(w http.ResponseWriter, r *http.Request) *AppError {
	var forgetPasswordParams ForgetPasswordParams
	if err := json.NewDecoder(r.Body).Decode(&forgetPasswordParams); err != nil {
		return DecodingRequestBodyError(err)
	}

	user, err := db.FindUserByEmail(forgetPasswordParams.Email)
	if err != nil {
		return DatabaseError(err, fmt.Sprintf("FindUserByEmail : %s", forgetPasswordParams.Email))
	}

	if user == nil {
		return BadRequestError(err, "No user with that email")
	}

	db.DeleteResetPassword(forgetPasswordParams.Email)

	token, err := uuid.NewUUID()
	if err != nil {
		return InternalError(err, "uuid.NewUUID()")
	}

	err = db.SaveResetPasswordToken(forgetPasswordParams.Email, token.String())
	if err != nil {
		return DatabaseError(err, fmt.Sprintf("SaveResetToken : %s - %s", forgetPasswordParams.Email, token))
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(true); err != nil {
		return EncodingResponseError(err, true)
	}

	return nil
}

func ResettingPassword(w http.ResponseWriter, r *http.Request) *AppError {
	return nil
}
