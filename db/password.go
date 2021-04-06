package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"guitou.cm/msvc/auth/models"
)

var pwdCollection = Database().Collection("passwords")

func SaveResetToken(email, token string) error {
	now := time.Now()
	data := models.ResetPassword{
		primitive.NilObjectID,
		email,
		token,
		now,
		now.Add(time.Hour * time.Duration(10)),
	}

	_, err := pwdCollection.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}
