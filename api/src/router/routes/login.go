package routes

import (
	"api/src/controllers"
	"net/http"
)

var routeLogin = Rota{
	URI:                "/login",
	Metodo:             http.MethodPost,
	Funcao:             controllers.Login,
	RequerAutenticacao: false,
}
