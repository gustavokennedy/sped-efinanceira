package repositories

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"sped-efinanceira/models"
)

type AutenticarRepository struct {
	db *mongo.Database
}

func NovoAutenticarRepository(dbURL, dbName string) (*AutenticarRepository, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return &AutenticarRepository{db: db}, nil
}

func (ar *AutenticarRepository) BuscarUsuarioPorEmail(email string) (*models.Usuario, error) {
	filter := bson.M{"email": email}

	var user models.Usuario
	err := ar.db.Collection("usuarios").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}
