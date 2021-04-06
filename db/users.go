package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"guitou.cm/msvc/auth/models"
)

var usersCollection = Database().Collection("users")

func SaveUser(u models.User) (primitive.ObjectID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	log.Println("Hast: ", hash)
	if err != nil {
		return primitive.NilObjectID, err
	}

	u.Password = string(hash)
	insertResult, err := usersCollection.InsertOne(context.Background(), u)
	if err != nil {
		return primitive.NilObjectID, err
	}
	log.Println("models.User inserted : ", u)

	return insertResult.InsertedID.(primitive.ObjectID), nil
}

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := usersCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func UpdateUserPassword(u models.User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	updateResult, err := usersCollection.UpdateByID(
		context.Background(),
		u.ID,
		bson.M{
			"$set": bson.M{
				"password": string(hash),
			},
		},
	)
	if err != nil {
		log.Println("Update Error : ", err)
		return err
	}
	log.Println("Update : ", updateResult.ModifiedCount, updateResult.MatchedCount)

	return nil
}
