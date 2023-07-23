package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Users struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Fullname  string             `bson:"fullname"`
	Email     string             `bson:"email"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	Role      string             `bson:"role"`
	Timestamp time.Time          `bson:"timestamp"`	
}