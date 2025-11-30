package models

import (
	"time"
)

type User struct {
    ID        string    `json:"id" bson:"_id"`
    Name      string    `json:"name" bson:"name" binding:"required,min=2,max=100"`
    Email     string    `json:"email" bson:"email" binding:"required,email"`
    Password  string    `json:"password,omitempty" bson:"password" binding:"required,min=6"`
    Role      string    `json:"role" bson:"role" binding:"required,oneof=admin user"`
    Phone     string    `json:"phone,omitempty" bson:"phone,omitempty"`
    Avatar    string    `json:"avatar,omitempty" bson:"avatar,omitempty"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type UpdateUser struct {
    Name   string `json:"name,omitempty" bson:"name,omitempty" binding:"omitempty,min=2,max=100"`
    Email  string `json:"email,omitempty" bson:"email,omitempty" binding:"omitempty,email"`
    Phone  string `json:"phone,omitempty" bson:"phone,omitempty"`
    Avatar string `json:"avatar,omitempty" bson:"avatar,omitempty"`
}

type RegisterInput struct {
    Name     string `json:"name" binding:"required,min=2,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Phone    string `json:"phone,omitempty"`
}

type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}
