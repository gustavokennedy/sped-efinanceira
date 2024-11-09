package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usuario struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Nome      string             `json:"nome" bson:"nome"`
	Email     string             `json:"email" bson:"email"`
	Senha     string             `json:"senha" bson:"senha" `
	Documento string             `json:"documento" bson:"documento"`
	Telefone  string             `json:"telefone" bson:"telefone"`
	Cidade    string             `json:"cidade" bson:"cidade"`
	PerfilID  string             `json:"perfil_id" bson:"perfil_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt time.Time          `json:"deleted_at" bson:"deleted_at"`
}
