package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Função para estabelecer a conexão com o banco de dados
func Connect(dbURL, dbName string) (*mongo.Client, *mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}

	db := client.Database(dbName)
	log.Println("✅ Conexão com o banco de dados estabelecida!")

	return client, db, nil
}

// Função para verificar o status da conexão com o banco de dados
func CheckConnection(client *mongo.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Realiza um ping no banco de dados para verificar a conexão
	err := client.Ping(ctx, nil)
	return err == nil
}
