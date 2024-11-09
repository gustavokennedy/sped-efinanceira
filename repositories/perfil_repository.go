package repositories

import (
	"context"
	"fmt"
	"log"
	"sped-efinanceira/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type PerfilRepositorio struct {
	db *mongo.Database
}

// Criar Perfil

func NovoPerfilRepositorio(dbURL, dbName string) (*PerfilRepositorio, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return &PerfilRepositorio{db: db}, nil
}

func (ur *PerfilRepositorio) CriarPerfil(perfil *models.Perfil) (*models.Perfil, error) {
	perfil.ID = primitive.NewObjectID()

	_, err := ur.db.Collection("perfis").InsertOne(context.Background(), perfil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return perfil, nil
}

// Listar todos Perfis
func (ur *PerfilRepositorio) ListarTodosPerfis() ([]models.Perfil, error) {
	var perfis []models.Perfil

	cur, err := ur.db.Collection("perfis").Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var perfil models.Perfil
		err := cur.Decode(&perfil)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		perfis = append(perfis, perfil)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return perfis, nil
}

// Listar Perfil por ID

func (ur *PerfilRepositorio) ListarPerfilPorID(id string) (*models.Perfil, error) {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var perfil models.Perfil
	err = ur.db.Collection("perfis").FindOne(context.Background(), filter).Decode(&perfil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &perfil, nil
}

// Buscar Nome
func (r *PerfilRepositorio) BuscarPerfilPorNome(nome string) (*models.Perfil, error) {
	collection := r.db.Collection("perfis")

	filter := bson.M{"nome": nome}

	var perfil models.Perfil
	err := collection.FindOne(context.Background(), filter).Decode(&perfil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Perfil não encontrado
		}
		return nil, err
	}

	return &perfil, nil
}

// Editar
func (ur *PerfilRepositorio) EditarPerfil(perfil *models.Perfil) error {
	filter := bson.M{"_id": perfil.ID}

	update := bson.M{
		"$set": bson.M{
			"descricao":  perfil.Descricao,
			"updated_at": time.Now(),
		},
	}

	_, err := ur.db.Collection("perfis").UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Perfil editado com sucesso!")
	return nil
}

// Deletar
func (ur *PerfilRepositorio) DeletarPerfil(id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println(err)
		return err
	}

	filter := bson.M{"_id": objectID}

	// Verificar se o Perfil existe
	count, err := ur.db.Collection("perfis").CountDocuments(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return err
	}

	if count == 0 {
		return fmt.Errorf("Perfil não encontrado!")
	}

	// Deletar o Perfil
	_, err = ur.db.Collection("perfis").DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Perfil deletado com sucesso!")
	return nil
}
