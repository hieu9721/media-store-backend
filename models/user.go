package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Name      string             `json:"name" bson:"name" binding:"required,min=2,max=100"`
    Email     string             `json:"email" bson:"email" binding:"required,email"`
    Age       int                `json:"age" bson:"age" binding:"required,gte=0,lte=150"`
    Phone     string             `json:"phone,omitempty" bson:"phone,omitempty"`
    Address   string             `json:"address,omitempty" bson:"address,omitempty"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type UpdateUser struct {
    Name    string `json:"name,omitempty" bson:"name,omitempty" binding:"omitempty,min=2,max=100"`
    Email   string `json:"email,omitempty" bson:"email,omitempty" binding:"omitempty,email"`
    Age     int    `json:"age,omitempty" bson:"age,omitempty" binding:"omitempty,gte=0,lte=150"`
    Phone   string `json:"phone,omitempty" bson:"phone,omitempty"`
    Address string `json:"address,omitempty" bson:"address,omitempty"`
}
