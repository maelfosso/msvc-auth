package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	log.Println("CheckPassword Into : ", err)
	return err == nil
	// if err != nil {
	// 	return false, err
	// }

	// return true, nil
}
