package seeders

import (
	"log"
	"os"
	"sped-efinanceira/models"
	"sped-efinanceira/repositories"
)

func ConfiguraRepositorios() repositories.PerfilRepositorio {

	dbURL := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

	perfilRepo, err := ConfiguraPerfilRepo(dbURL, dbName)
	if err != nil {
		log.Fatal("Erro com Repo Perfil:", err)
	}

	return *perfilRepo
}

func ConfiguraPerfilRepo(dbURL string, dbName string) (*repositories.PerfilRepositorio, error) {
	return repositories.NovoPerfilRepositorio(dbURL, dbName)
}

// Perfil
func SeedPerfis(perfilRepo *repositories.PerfilRepositorio) {
	perfisExistentes, err := perfilRepo.ListarTodosPerfis()
	if err != nil {
		log.Println("Erro ao listar Perfis:", err)
		return
	}

	perfisMap := make(map[string]bool)
	for _, perfil := range perfisExistentes {
		perfisMap[perfil.Nome] = true
	}

	// Criar o perfil de Administrador, se não existir
	if !perfisMap["Admin"] {
		perfilAdmin := &models.Perfil{
			Nome:      "Admin",
			Descricao: "Permissões de administrador.",
		}

		_, err = perfilRepo.CriarPerfil(perfilAdmin)
		if err != nil {
			log.Println("Erro ao criar Perfil:", err)
			return
		}

		log.Printf("🌱 Seed: Perfil '%s' criado com sucesso!", perfilAdmin.Nome)
	} else {
		log.Println("🌱 Seed: Perfil 'Administrador' já existe.")
	}

	// Criar o perfil de Clientes, se não existir
	if !perfisMap["Clientes"] {
		perfilClientes := &models.Perfil{
			Nome:      "Clientes",
			Descricao: "Permissões de clientes.",
		}

		_, err = perfilRepo.CriarPerfil(perfilClientes)
		if err != nil {
			log.Println("Erro ao criar Perfil:", err)
			return
		}

		log.Printf("🌱 Seed: Perfil '%s' criado com sucesso!", perfilClientes.Nome)
	} else {
		log.Println("🌱 Seed: Perfil 'Clientes' já existe.")
	}
}
