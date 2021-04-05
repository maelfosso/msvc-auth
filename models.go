package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SignInParams struct {
	Email    string `json:"email"`
	Password string `json:"email"`
}

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FullName string             `json:"fullname" bson:"fullname"`
	Email    string             `json:"email" bson:"email"`
	Phone    string             `json:"phone" bson:"phone"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

var collection = Database().Collection("users")

func (u *User) Save() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	log.Println("Hast: ", hash)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	insertResult, err := collection.InsertOne(context.Background(), u)
	if err != nil {
		return err
	}
	log.Println("User inserted : ", u)

	u.ID = insertResult.InsertedID.(primitive.ObjectID)
	return nil
}

// func (u *User) MarshallJSON() ([]byte, error) {
// 	type Alias User

// 	return json.Marshal(&struct {
// 		*Alias
// 		Password *string `json:"password"`
// 		Cost     *int    `json:"cost"`
// 	}{
// 		Password: nil,
// 		Cost:     nil,
// 		Alias:    (*Alias)(u),
// 	})
// }

func FindUserByEmail(email string) (*User, error) {
	var user User
	if err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func CheckPassword(hash, password string, cost int) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}