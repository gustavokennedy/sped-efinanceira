package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Perfil struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Nome  string `json:"nome"`
	Descricao string             `json:"descricao" bson:"descricao"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt time.Time          `json:"deleted_at" bson:"deleted_at"`
}