package router

import (
	"api/src/router/routes"

	"github.com/gorilla/mux"
)

// Gerar vai retornar um router
func Gerar() *mux.Router {
	r := mux.NewRouter()
	return routes.Config(r)
}
