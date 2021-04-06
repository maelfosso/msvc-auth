package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"guitou.cm/msvc/auth/models"
)

var pwdCollection = Database().Collection("passwords")

func SaveResetPasswordToken(email, token string) error {
	now := time.Now()
	data := models.ResetPassword{
		primitive.NilObjectID,
		email,
		token,
		now.Add(time.Minute * 10),
		now,
	}

	_, err := pwdCollection.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}

func DeleteResetPassword(email string) (*models.ResetPassword, error) {
	var resetPwd models.ResetPassword

	if err := pwdCollection.FindOneAndDelete(context.Background(), bson.M{"email": email}).Decode(&resetPwd); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &resetPwd, nil
}

func FindResetPassword(email, token string) (*models.ResetPassword, error) {
	var resetPwd models.ResetPassword

	if err := pwdCollection.FindOne(context.Background(), bson.M{"email": email, "token": token}).Decode(&resetPwd); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &resetPwd, nil
}
