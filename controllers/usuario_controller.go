package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"sped-efinanceira/common"
	"sped-efinanceira/models"
	"sped-efinanceira/repositories"
)

type RespostaUsuario struct {
	Usuario *Usuario `json:"usuario"`
}
type Usuario struct {
	*models.Usuario
	Perfil *models.Perfil `json:"perfil"`
}

type UsuarioController struct {
	repo       *repositories.UsuarioRepositorio
	perfilRepo *repositories.PerfilRepositorio
	authRepo   *repositories.AutenticarRepository
}

func NovoUsuarioController(repo *repositories.UsuarioRepositorio, perfilRepo *repositories.PerfilRepositorio, authRepo *repositories.AutenticarRepository) *UsuarioController {
	return &UsuarioController{
		repo:       repo,
		perfilRepo: perfilRepo,
		authRepo:   authRepo,
	}
}

// generateJWT gera um token JWT para um usuário com base em seu ID
func (uc *UsuarioController) generateJWT(userID string) (string, error) {

	objectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return "", err
	}

	// Define as claims do token (informações sobre o usuário)
	claims := jwt.MapClaims{
		"sub": objectID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expira em 24 horas
	}

	// Cria um novo token com as claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assina o token com a chave secreta
	signedToken, err := token.SignedString([]byte("oc-energia"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Criar Usuário
func (uc *UsuarioController) CriarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario models.Usuario
	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Pedido inválido!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	validate := validator.New()
	if err := validate.Struct(usuario); err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Campos inválidos!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	err = uc.repo.CriarUsuario(&usuario, usuario.PerfilID)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao criar Usuário!",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuario)

}

// Autenticar Usuário
func (uc *UsuarioController) AutenticarUsuario(w http.ResponseWriter, r *http.Request) {
	var authData models.AuthData
	err := json.NewDecoder(r.Body).Decode(&authData)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Pedido inválido!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	// Vai verificar se usuário exisste
	user, err := uc.authRepo.BuscarUsuarioPorEmail(authData.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			RespostaComErro := common.RespostaComErro{
				Error:   "Usuário não encontrado!",
				Message: err.Error(),
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(RespostaComErro)
		} else {
			log.Println(err)
			RespostaComErro := common.RespostaComErro{
				Error:   "Falhha ao autenticar Usuário!",
				Message: err.Error(),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(RespostaComErro)
		}
		return
	}

	// Valida a senha
	err = bcrypt.CompareHashAndPassword([]byte(user.Senha), []byte(authData.Senha))
	if err != nil {
		RespostaComErro := common.RespostaComErro{
			Error:   "Credenciais inválidas!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	// Gera Token
	token, err := uc.generateJWT(user.ID.Hex())
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao gerar token de autenticação!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	Resposta := models.AuthResponse{
		Token: token,
	}
	log.Printf("Usuário logado: %s", user.Nome)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Resposta)
}

// Usuario Logado

// Obter informações do usuário logado
func (uc *UsuarioController) ObterInformacoesUsuarioLogado(w http.ResponseWriter, r *http.Request) {
	// Extrair o token do cabeçalho de autorização
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// Validar e decodificar o token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar o método de assinatura do token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de assinatura inválido")
		}
		// Retornar a chave secreta usada para assinar o token
		return []byte("efinanceira"), nil
	})

	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Token inválido!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	// Verificar se o token é válido e não expirou
	if !token.Valid {
		RespostaComErro := common.RespostaComErro{
			Error:   "Token inválido ou expirado!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	// Extrair o ID do usuário do token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		RespostaComErro := common.RespostaComErro{
			Error:   "Token inválido!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		RespostaComErro := common.RespostaComErro{
			Error:   "Token inválido!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	// Buscar o usuário no banco de dados
	usuario, err := uc.repo.ListarUsuarioPorID(userID)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao buscar Usuário!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	perfil, err := uc.perfilRepo.ListarPerfilPorID(usuario.PerfilID)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao buscar Perfil do Usuário!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao buscar Controladores do Usuário!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	Usuario := &Usuario{
		Usuario: usuario,
		Perfil:  perfil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Usuario)
}

// Listar por ID
func (uc *UsuarioController) ListarUsuarioPorID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	usuario, err := uc.repo.ListarUsuarioPorID(id)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao buscar Usuário!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	perfil, err := uc.perfilRepo.ListarPerfilPorID(usuario.PerfilID)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao buscar Perfil do Usuário!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao buscar Controladores do Usuário!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	Usuario := &Usuario{
		Usuario: usuario,
		Perfil:  perfil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Usuario)
}

// Listar Usuários
func (uc *UsuarioController) BuscarPerfil(perfilID string) (*models.Perfil, error) {
	perfil, err := uc.perfilRepo.ListarPerfilPorID(perfilID)
	if err != nil {
		return nil, err
	}
	return perfil, nil
}

// ListarUsuarios lista todos os usuários
func (uc *UsuarioController) ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := uc.repo.ListarUsuarios()
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao buscar Usuários!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	usuariosComPerfil := make([]*Usuario, 0, len(usuarios))

	for _, usuario := range usuarios {

		perfil, err := uc.perfilRepo.ListarPerfilPorID(usuario.PerfilID)
		if err != nil {
			log.Println(err)
			RespostaComErro := common.RespostaComErro{
				Error:   "Falha ao buscar Perfil do Usuário!",
				Message: err.Error(),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(RespostaComErro)
			return
		}

		Usuario := &Usuario{
			Usuario: usuario,
			Perfil:  perfil,
		}

		usuariosComPerfil = append(usuariosComPerfil, Usuario)
	}

	// Adicionando a contagem de usuários à resposta JSON
	resposta := struct {
		TotalUsuarios int        `json:"total_usuarios"`
		Usuarios      []*Usuario `json:"usuarios"`
	}{
		TotalUsuarios: len(usuarios),
		Usuarios:      usuariosComPerfil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resposta)
}

// Atualizar Usuário

func (uc *UsuarioController) AtualizarUsuario(w http.ResponseWriter, r *http.Request) {

	var usuario models.Usuario
	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Pedido inválido!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	validate := validator.New()
	if err := validate.Struct(usuario); err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Campos inválidos!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	// Verificar se o usuário existe antes de atualizá-lo
	vars := mux.Vars(r)
	id := vars["id"]

	existeUsuario, err := uc.repo.ListarUsuarioPorID(id)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Usuário não encontrado!",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	if existeUsuario == nil {
		RespostaComErro := common.RespostaComErro{
			Error:   "Controlador não encontrado!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	err = uc.repo.AtualizarUsuario(existeUsuario.ID, &usuario)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao editar Usuário!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	Resposta := usuario

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Resposta)
}

// Deletar Usuário
func (uc *UsuarioController) DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	err := uc.repo.DeletarUsuario(id)
	if err != nil {
		if err.Error() == "Usuário não encontrado" {
			// Usuário não encontrado, retornar mensagem em JSON
			RespostaComErro := common.RespostaComErro{
				Error:   "Usuário não encontrado!",
				Message: err.Error(),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(RespostaComErro)
			return
		}

		log.Println(err)
		http.Error(w, "Falha ao deletar usuário!", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
