package models

import (
	"api/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

//Representa o modelo do usuário
type User struct {
	ID        uint64    `json:"id,omitempty"`
	Nome      string    `json:"nome,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"CreatedAt,omitempty"`
}

//Chama os métodos para validar e formatar o usuário
func (usuario *User) Preparar(etapa string) error {
	if erro := usuario.validar(etapa); erro != nil {
		return erro
	}
	if erro := usuario.formatar(etapa); erro != nil {
		return erro
	}
	return nil
}

func (usuario *User) validar(etapa string) error {
	if usuario.Nome == "" {
		return errors.New("O nome é obrigatório e não pode ser vazio")
	}

	if usuario.Nick == "" {
		return errors.New("O Nick é obrigatório e não pode ser vazio")
	}
	if usuario.Email == "" {
		return errors.New("O Email é obrigatório e não pode ser vazio")
	}
	if etapa == "cadastro" && usuario.Password == "" {
		return errors.New("A senha é obrigatória e não pode ser vazia")
	}
	if erro := checkmail.ValidateFormat(usuario.Email); erro != nil {
		return errors.New("O e-mail inserido é inválido")
	}
	return nil
}

func (usuario *User) formatar(etapa string) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)

	if etapa == "cadastro" {
		passwordHash, erro := security.Hash(usuario.Password)
		if erro != nil {
			return erro
		}

		usuario.Password = string(passwordHash)
	}
	return nil
}
