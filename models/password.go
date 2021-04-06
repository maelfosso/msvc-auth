package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResetPassword struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Token     string             `json:"token" bson:"token"`
	ExpiredAt time.Time          `json:"expiredAt,omitempty" bson:"expiredAt"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt"`
}
