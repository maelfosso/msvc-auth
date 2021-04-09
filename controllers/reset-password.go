package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"guitou.cm/msvc/auth/db"
	"guitou.cm/msvc/auth/rabbitmq"
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

	data, err := db.SaveResetPasswordToken(forgetPasswordParams.Email, token.String())
	if err != nil {
		return DatabaseError(err, fmt.Sprintf("SaveResetToken : %s - %s", forgetPasswordParams.Email, token))
	}

	if b, err := json.Marshal(data); err == nil {
		log.Println("Marshelling before Publishing : %v", b)
		rabbitmq.Publish(b, "auth.password.forget")
	} else {
		log.Println("[rset-password] Error occurred when marshelling data")
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(true); err != nil {
		return EncodingResponseError(err, true)
	}

	return nil
}

func VerifyResetPasswordToken(w http.ResponseWriter, r *http.Request) *AppError {
	var params VerifyResetPasswordTokenParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return DecodingRequestBodyError(err)
	}

	resetToken, err := db.FindResetPassword(params.Email, params.Token)
	if err != nil {
		return DatabaseError(err, fmt.Sprintf("FindResetPassword : %s - %s", params.Email, params.Token))
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resetToken != nil); err != nil {
		return EncodingResponseError(err, resetToken != nil)
	}

	return nil
}

func ResettingPassword(w http.ResponseWriter, r *http.Request) *AppError {
	var params ResettingPasswordParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return DecodingRequestBodyError(err)
	}

	if params.Password != params.ConfirmPassword {
		return BadRequestError(fmt.Errorf("password don't match"), "Password don't match")
	}

	_, err := db.FindResetPassword(params.Email, params.Token)
	if err != nil {
		return DatabaseError(err, fmt.Sprintf("FindResetPassword : %s - %s", params.Email, params.Token))
	}

	user, err := db.FindUserByEmail(params.Email)
	if err != nil {
		return DatabaseError(err, "Impossible to retreive the user")
	}

	if user == nil {
		return BadRequestError(err, "No user with that email")
	}

	err = db.UpdateUserPassword(*user, params.Password)
	if err != nil {
		return DatabaseError(err, "Error when updating the password")
	}

	// Delete the token document

	if b, err := json.Marshal(user); err == nil {
		log.Println("Marshelling before Publishing : %v", b)
		rabbitmq.Publish(b, "auth.password.reset")
	} else {
		log.Println("[rset-password] Error occurred when marshelling data")
	}

	return nil
}
