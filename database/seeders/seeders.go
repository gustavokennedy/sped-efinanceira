package seeders

import (
	"log"
	"os"
	"sped-efinanceira/models"
	"sped-efinanceira/repositories"
)

func ConfiguraRepositorios() (repositories.UsuarioRepositorio, repositories.PerfilRepositorio) {

	dbURL := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

	perfilRepo, err := ConfiguraPerfilRepo(dbURL, dbName)
	if err != nil {
		log.Fatal("Erro com Repo Perfil:", err)
	}

	usuarioRepo, err := ConfiguraUsuarioRepo(dbURL, dbName)
	if err != nil {
		log.Fatal("Erro com Repo Usuário:", err)
	}

	return *usuarioRepo, *perfilRepo
}

func ConfiguraUsuarioRepo(dbURL string, dbName string) (*repositories.UsuarioRepositorio, error) {
	return repositories.NovoUsuarioRepository(dbURL, dbName)
}

func ConfiguraPerfilRepo(dbURL string, dbName string) (*repositories.PerfilRepositorio, error) {
	return repositories.NovoPerfilRepositorio(dbURL, dbName)
}

// Usuários
func SeedUsuarios(usuarioRepo *repositories.UsuarioRepositorio, perfilRepo *repositories.PerfilRepositorio) {
	SeedPerfis(perfilRepo)

	usuarios, err := usuarioRepo.ListarUsuarios()
	if err != nil {
		log.Println("Erro ao listar Usuários:", err)
		return
	}

	if len(usuarios) == 0 {
		user := &models.Usuario{
			Nome:      "Gustavo",
			Email:     "gustavo@overall.cloud",
			Senha:     "teste123",
			Documento: "89613078940",
			Telefone:  "4737540330",
			Cidade:    "Aurora",
		}

		// Verifica se o usuário já existe pelo email
		existingUser, err := usuarioRepo.BuscarUsuarioPorEmail(user.Email)
		if err != nil {
			log.Println("Erro ao buscar Usuário por email:", err)
			return
		}

		if existingUser != nil {
			// Logar que o usuário já existe e retornar
			log.Println("🌱 Seed: Usuário Gustavo já existe.")
			return
		}

		// Verifica o perfil de administrador
		perfilAdministrador, err := perfilRepo.BuscarPerfilPorNome("Admin")
		if err != nil {
			log.Println("Erro ao buscar Perfil:", err)
			return
		}

		if perfilAdministrador == nil {
			log.Println("Perfil 'Admin' não encontrado.")
			return
		}

		user.PerfilID = perfilAdministrador.ID.Hex()

		err = usuarioRepo.CriarUsuario(user, user.PerfilID)
		if err != nil {
			log.Println("Erro ao criar Usuário:", err)
			return
		}

		log.Printf("🌱 Seed: Usuário '%s' criado com sucesso!", user.Nome)
	} else {
		// Caso o usuário já exista, pode logar isso aqui também
		log.Println("🌱 Seed: Usuários já existem. Nenhum novo usuário criado.")
	}
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
