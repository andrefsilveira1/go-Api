package middleware

import (
	"api/src/autenticar"
	"api/src/respostas"
	"fmt"
	"net/http"
)

// Imprime no terminal
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n $s $s $s", r.Method, r.RequestURI, r.Host)
		next(w, r)
	}
}

// Middleware para verificar se o usuário está autenticado
func Autenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := autenticar.ValidateToken(r); erro != nil {
			respostas.Erro(w, http.StatusUnauthorized, erro)
			return
		}
		fmt.Println("Validando...")
		next(w, r)
	}
}
