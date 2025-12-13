package models

import (
	"time"

	"go.mondgo.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectId `bson:"_id"`
	First_name    *string            `json:"first_name" validata:"required,min=2,max=100"`
	Last_name     *string            `json:"last_name" validata:"required,min=2,max=100"`
	Password      *string            `json:"password" validata:"required,min=6"`
	Email         *string            `json:"email" validate:"email"`
	Phone         *string            `json:"phone" validate:"required"`
	Token         *string            `json:"token"`
	User_type     *string            `json:"user_type" validate:"required, eq=ADMIN | eq=USER"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}
