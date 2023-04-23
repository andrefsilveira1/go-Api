package controllers

import (
	"api/src/autenticar"
	"api/src/data"
	"api/src/models"
	"api/src/repositories"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Cria
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	bodyRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	var user models.User
	if erro = json.Unmarshal(bodyRequest, &user); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	if erro = user.Preparar("cadastro"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := data.Connect()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewRepository(db)
	ID, erro := repositorio.CriarUsuario(user)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusCreated, user)
	w.Write([]byte(fmt.Sprintf("ID do usuário %d", ID)))
}

// Busca todos
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nome := strings.ToLower(r.URL.Query().Get("usuario"))

	db, erro := data.Connect()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewRepository(db)
	usuarios, erro := repositorio.Buscar(nome)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, usuarios)
}

// Busca um
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	usuarioId, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	db, erro := data.Connect()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewRepository(db)
	usuario, erro := repositorio.BuscarPorId(usuarioId)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, usuario)
	w.Write([]byte("Buscando usuário"))
}

// Atualiza um
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	userIdToken, erro := autenticar.GetUserId(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
	}
	fmt.Print(userIdToken)

	if usuarioID != userIdToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Você não tem autorização para essa requisição"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	var usuario models.User
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return

	}
	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := data.Connect()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewRepository(db)
	if erro = repositorio.Atualizar(usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return

	}
	respostas.JSON(w, http.StatusNoContent, nil)
	w.Write([]byte("Atualizando usuário"))
}

// Deleta um
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	userIdToken, erro := autenticar.GetUserId(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
	}

	if usuarioID != userIdToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Você não tem autorização para essa requisição"))
	}
	db, erro := data.Connect()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewRepository(db)
	if erro = repositorio.Deletar(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
	w.Write([]byte("Deletando usuário"))
}

// Seguir usuário
func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerId, erro := autenticar.GetUserId(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	userId, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if followerId == userId {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	db, erro := data.Connect()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewRepository(db)
	if erro = repositorio.Follow(userId, followerId); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// Função criada para permitir que um usuário pare de seguir outro usuário
func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followerId, erro := autenticar.GetUserId(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	userId, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if followerId == userId {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível parar de seguir seu próprio usuário"))
		return
	}
	db, erro := data.Connect()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositories := repositories.NewRepository(db)
	if erro = repositories.UnfollowUser(userId, followerId); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
