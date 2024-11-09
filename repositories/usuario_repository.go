package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"

	"sped-efinanceira/models"
)

type UsuarioRepositorio struct {
	db *mongo.Database
}

// Criar Usuário

func NovoUsuarioRepository(dbURL, dbName string) (*UsuarioRepositorio, error) {
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
	return &UsuarioRepositorio{db: db}, nil
}

func (ur *UsuarioRepositorio) CriarUsuario(usuario *models.Usuario, perfilID string) error {

	usuario.ID = primitive.NewObjectID()
	// Gerar o hash da senha
	senhaEncriptada, err := bcrypt.GenerateFromPassword([]byte(usuario.Senha), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}

	// Substituir a senha pelo hash gerado
	usuario.Senha = string(senhaEncriptada)

	// Atribuir o PerfilID ao usuário
	usuario.PerfilID = perfilID

	_, err = ur.db.Collection("usuarios").InsertOne(context.Background(), usuario)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Usuário criado com sucesso!")
	return nil
}

// Listar

func (sr *UsuarioRepositorio) ListarUsuarios() ([]*models.Usuario, error) {
	filter := bson.M{} // Filtro vazio para listar todos os usuários

	cursor, err := sr.db.Collection("usuarios").Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var usuarios []*models.Usuario
	for cursor.Next(context.Background()) {
		var usuario models.Usuario
		err := cursor.Decode(&usuario)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		usuarios = append(usuarios, &usuario)
	}

	if err := cursor.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return usuarios, nil
}

// Listar Usuário por ID
func (ur *UsuarioRepositorio) ListarUsuarioPorID(id string) (*models.Usuario, error) {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var usuario models.Usuario
	err = ur.db.Collection("usuarios").FindOne(context.Background(), filter).Decode(&usuario)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Usuário listado com sucesso!")
	return &usuario, nil
}

// Buscar por Email
func (ur *UsuarioRepositorio) BuscarUsuarioPorEmail(email string) (*models.Usuario, error) {
	filter := bson.M{"email": email}

	var usuario models.Usuario
	err := ur.db.Collection("usuarios").FindOne(context.Background(), filter).Decode(&usuario)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Usuário não encontrado
		}
		log.Println(err)
		return nil, err
	}

	return &usuario, nil
}

// Atualizar Usuário

func (ur *UsuarioRepositorio) AtualizarUsuario(id primitive.ObjectID, usuario *models.Usuario) error {
	usuario.ID = id
	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{
			"nome":       usuario.Nome,
			"email":      usuario.Email,
			"documento":  usuario.Documento,
			"telefone":   usuario.Telefone,
			"cidade":     usuario.Cidade,
			"perfil_id":  usuario.PerfilID,
			"updated_at": time.Now(),
		},
	}

	if usuario.Senha != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usuario.Senha), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			return err
		}
		update["$set"].(bson.M)["senha"] = string(hashedPassword)
	}

	_, err := ur.db.Collection("usuarios").UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Usuário editado com sucesso!")
	return nil
}

// Deletar Usuário
func (ur *UsuarioRepositorio) DeletarUsuario(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return err
	}

	filter := bson.M{"_id": objectID}

	// Verificar se o usuário existe
	count, err := ur.db.Collection("usuarios").CountDocuments(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return err
	}

	if count == 0 {
		// Usuário não encontrado, retornar um erro ou uma mensagem adequada
		return fmt.Errorf("Usuário não encontrado!")
	}

	// Deletar o usuário
	_, err = ur.db.Collection("usuarios").DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Usuário deletado com sucesso!")
	return nil
}
