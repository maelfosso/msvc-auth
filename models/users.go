package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FullName string             `json:"fullname" bson:"fullname"`
	Email    string             `json:"email" bson:"email"`
	Phone    string             `json:"phone" bson:"phone"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User

	u.Password = ""
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&u),
	})
}
