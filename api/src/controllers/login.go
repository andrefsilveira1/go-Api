package controllers

import (
	"api/src/autenticar"
	"api/src/data"
	"api/src/models"
	"api/src/repositories"
	"api/src/respostas"
	"api/src/security"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.User
	if erro = json.Unmarshal(body, &usuario); erro != nil {
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
	usuarioTemp, erro := repositorio.FindByEmail(usuario.Email)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	if erro = security.Verify(usuario.Password, usuarioTemp.Password); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	w.Write([]byte("Você está logado!"))
	fmt.Println(usuarioTemp)
	token, _ := autenticar.GenerateToken(usuarioTemp.ID)
	fmt.Println(token)
	w.Write([]byte(token))
}
