package autenticar

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//Gerar token JWT com permissão do usuário
func GenerateToken(usuarioId uint64) (string, error) {
	permission := jwt.MapClaims{}
	permission["authorized"] = true
	permission["exp"] = time.Now().Add(time.Hour * 4).Unix()
	permission["usuarioId"] = usuarioId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permission)
	return token.SignedString([]byte(config.SecretKey)) //Secret
}

// Função para autenticar o token passado
func ValidateToken(r *http.Request) error {
	tokenTemp := extractToken(r)
	token, erro := jwt.Parse(tokenTemp, returnVerifyKey)
	if erro != nil {
		return erro
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return errors.New("Token inválido")
}

// Pega o Id do usuário passado
func GetUserId(r *http.Request) (uint64, error) {
	tokenTemp := extractToken(r)
	token, erro := jwt.Parse(tokenTemp, returnVerifyKey)
	if erro != nil {
		return 0, erro
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usuarioId, erro := strconv.ParseUint(fmt.Sprintf("%.0f", permissions["usuarioId"]), 10, 64)
		if erro != nil {
			return 0, erro
		}

		return usuarioId, nil
	}
	return 0, nil
}
func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func returnVerifyKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("O método não é o esperado %v", token.Header["alg"])
	}
	return config.SecretKey, nil
}
