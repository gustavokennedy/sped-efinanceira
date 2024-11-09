package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"

	"sped-efinanceira/common"
	"sped-efinanceira/models"
	"sped-efinanceira/repositories"
)

type PerfilController struct {
	repo *repositories.PerfilRepositorio
}

func NovoPerfilController(repo *repositories.PerfilRepositorio) *PerfilController {
	return &PerfilController{repo: repo}
}

// Criar Perfil
func (uc *PerfilController) CriarPerfil(w http.ResponseWriter, r *http.Request) {
	var perfil models.Perfil
	err := json.NewDecoder(r.Body).Decode(&perfil)
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

	// Validar o modelo Perfil
	validate := validator.New()
	if err := validate.Struct(perfil); err != nil {
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

	perfilCriado, err := uc.repo.CriarPerfil(&perfil)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao criar Perfil!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(perfilCriado)
}

// Listar

// BuscarPerfilPorNome busca um perfil pelo nome
func (uc *PerfilController) BuscarPerfilPorNome(w http.ResponseWriter, r *http.Request) {
	// Para query parameter:
	// nome := r.URL.Query().Get("nome")

	// Para path parameter:
	nome := mux.Vars(r)["nome"]

	if nome == "" {
		RespostaComErro := common.RespostaComErro{
			Error:   "Nome não especificado",
			Message: "O parâmetro 'nome' é obrigatório.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	log.Printf("Buscando perfil com nome: %s", nome)

	perfil, err := uc.repo.BuscarPerfilPorNome(nome)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao buscar Perfil",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	if perfil == nil {
		RespostaComErro := common.RespostaComErro{
			Error:   "Perfil não encontrado",
			Message: "Nenhum perfil encontrado com o nome especificado.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perfil)
}

// Listar todos Perfis
func (uc *PerfilController) ListarTodosPerfis(w http.ResponseWriter, r *http.Request) {
	perfis, err := uc.repo.ListarTodosPerfis()
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao listar Perfis!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perfis)
}

func (uc *PerfilController) ListarPerfilPorID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	perfil, err := uc.repo.ListarPerfilPorID(id)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao listar Perfil!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perfil)
}

// Atualizar

func (uc *PerfilController) EditarPerfil(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var perfil models.Perfil
	err := json.NewDecoder(r.Body).Decode(&perfil)
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

	// Validar o modelo
	validate := validator.New()
	if err := validate.Struct(perfil); err != nil {
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
	existingPerfil, err := uc.repo.ListarPerfilPorID(id)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Perfil não encontrado!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	// Atualizar apenas os campos relevantes do perfil
	existingPerfil.Descricao = perfil.Descricao
	existingPerfil.UpdatedAt = time.Now()

	err = uc.repo.EditarPerfil(existingPerfil)
	if err != nil {
		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao atualizar Perfil!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	Resposta := common.Resposta{
		Message: "Perfil atualizado com sucesso!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Resposta)
}

// Deletar
func (uc *PerfilController) DeletarPerfil(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	err := uc.repo.DeletarPerfil(id)
	if err != nil {
		if err.Error() == "Perfil não encontrado" {
			RespostaComErro := common.RespostaComErro{
				Error:   "Perfil não encontrado!",
				Message: err.Error(),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(RespostaComErro)
			return
		}

		log.Println(err)
		RespostaComErro := common.RespostaComErro{
			Error:   "Falha ao deletar Perfil!",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaComErro)
		return
	}

	w.WriteHeader(http.StatusOK)
}
