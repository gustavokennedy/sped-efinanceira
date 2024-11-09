package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sped-efinanceira/controllers"
	"sped-efinanceira/database"
	"sped-efinanceira/middlewares"
	"sped-efinanceira/repositories"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Response struct {
	Message string `json:"message"`
}

type ApiStatus struct {
	AppStatus    string `json:"app_status"`
	DBConnection string `json:"db_connection"`
}

func EnviarEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Ler o corpo da requisição
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	// Decodificar o corpo da requisição JSON em uma estrutura de EmailRequest
	var emailReq EmailRequest
	err = json.Unmarshal(body, &emailReq)
	if err != nil {
		http.Error(w, "Erro ao decodificar o JSON da requisição", http.StatusBadRequest)
		return
	}

	// Criar uma instância do middleware de email
	emailMiddleware := middlewares.NovoEmailMiddleware()

	// Enviar o email
	err = emailMiddleware.SendEmail(emailReq.To, emailReq.Subject, emailReq.Body)
	if err != nil {
		log.Println("Erro ao enviar o email:", err)
		http.Error(w, "Erro ao enviar o email", http.StatusInternalServerError)
		return
	}

	// Responder com uma mensagem de sucesso
	response := Response{Message: "Email enviado com sucesso!"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Erro ao criar a resposta JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// ConfiguraRotas configura as rotas e recebe o client do MongoDB
func ConfiguraRotas(client *mongo.Client) *mux.Router {

	dbURL := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

	// Repositories
	usuarioRepo, err := repositories.NovoUsuarioRepository(dbURL, dbName)
	if err != nil {
		log.Fatal("Erro ao conectar ao repositório de usuários:", err)
	}

	perfilRepo, err := repositories.NovoPerfilRepositorio(dbURL, dbName)
	if err != nil {
		log.Fatal("Erro ao conectar ao repositório de perfis:", err)
	}

	authRepo, err := repositories.NovoAutenticarRepository(dbURL, dbName)
	if err != nil {
		log.Fatal("Erro ao conectar ao repositório de autenticação:", err)
	}

	// Inicializar o controlador de perfil
	perfilController := controllers.NovoPerfilController(perfilRepo)
	usuarioController := controllers.NovoUsuarioController(usuarioRepo, perfilRepo, authRepo)

	router := mux.NewRouter()

	// Rotas Publicas
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		dbStatus := "503"
		if database.CheckConnection(client) {
			dbStatus = "200"
		}

		response := ApiStatus{
			AppStatus:    "200",
			DBConnection: dbStatus,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	router.HandleFunc("/enviar-email", EnviarEmailHandler)
	router.HandleFunc("/logar", usuarioController.AutenticarUsuario).Methods("POST").Name("AutenticarUsuario")

	// Rotas Privadas

	// Grupo de rotas privadas com middleware de autenticação
	privateRoutes := router.PathPrefix("/").Subrouter()
	privateRoutes.Use(middlewares.AutenticarMiddleware)

	// Rotas para perfis
	privateRoutes.HandleFunc("/perfis", perfilController.CriarPerfil).Methods("POST").Name("CriarPerfil")
	privateRoutes.HandleFunc("/perfis", perfilController.ListarTodosPerfis).Methods("GET").Name("ListarPerfil")
	privateRoutes.HandleFunc("/perfis/{id}", perfilController.ListarPerfilPorID).Methods("GET").Name("ListarPerfil")
	privateRoutes.HandleFunc("/perfis/{id}", perfilController.EditarPerfil).Methods("PUT").Name("EditarPerfil")
	privateRoutes.HandleFunc("/perfis/{id}", perfilController.DeletarPerfil).Methods("DELETE").Name("DeletarPerfil")

	// Rotas para usuários
	privateRoutes.HandleFunc("/profile", usuarioController.ObterInformacoesUsuarioLogado).Methods("GET").Name("ListarUsuarioLogado")
	privateRoutes.HandleFunc("/usuarios", usuarioController.ListarUsuarios).Methods("GET").Name("ListarUsuarios")
	privateRoutes.HandleFunc("/usuarios/{id}", usuarioController.ListarUsuarioPorID).Methods("GET").Name("ListarUsuariosPorID")
	privateRoutes.HandleFunc("/usuarios", usuarioController.CriarUsuario).Methods("POST").Name("CriarUsuario")
	privateRoutes.HandleFunc("/usuarios/{id}", usuarioController.AtualizarUsuario).Methods("PUT").Name("AtualizarUsuario")
	privateRoutes.HandleFunc("/usuarios/{id}", usuarioController.DeletarUsuario).Methods("DELETE").Name("DeletarUsuario")

	return router
}
