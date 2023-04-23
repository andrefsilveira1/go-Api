package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

// Representa um repositório
type users struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *users {
	return &users{db}
}

func (u users) CriarUsuario(user models.User) (uint64, error) {
	statement, erro := u.db.Prepare("insert into usuarios (nome, nick, email, password) values (?,?,?,?)")
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(user.Nome, user.Nick, user.Email, user.Password)
	if erro != nil {
		return 0, erro
	}

	ID, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ID), nil
}

// Traz todos os usuários
func (repositorio users) Buscar(nome string) ([]models.User, error) {
	nome = fmt.Sprintf("%%%s%%", nome) //%nome%

	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, createdAt from usuarios where nome LIKE ? or nick LIKE ?", nome, nome,
	)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []models.User

	for linhas.Next() {
		var usuario models.User

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CreatedAt,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}
func (repositorio users) BuscarPorId(id uint64) (models.User, error) {
	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, createdAt from usuarios where id = ?", id,
	)
	if erro != nil {
		return models.User{}, erro
	}

	defer linhas.Close()

	var usuario models.User
	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CreatedAt,
		); erro != nil {
			return models.User{}, erro
		}
	}
	return usuario, nil
}

// Atualiza as informações de um usuário
func (repositorio users) Atualizar(ID uint64, usuario models.User) error {
	statement, erro := repositorio.db.Prepare("update usuarios set nome = ?, nick = ?, email = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, ID); erro != nil {
		return erro
	}

	return nil
}

// Deleta um usuário
func (repositorio users) Deletar(ID uint64) error {
	statement, erro := repositorio.db.Prepare("delete from usuarios where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}

// Busca usuário por email
func (repositorio users) FindByEmail(email string) (models.User, error) {
	linha, erro := repositorio.db.Query("select id, password from usuarios where email = ?", email)
	if erro != nil {
		return models.User{}, erro
	}
	defer linha.Close()

	var usuario models.User

	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Password); erro != nil {
			return models.User{}, erro
		}
	}
	return usuario, nil
}

// Adicionar usuário a ser seguido e um seguidor ao usuário seguido
func (repositorio users) Follow(userId, followerId uint64) error {
	statement, erro := repositorio.db.Prepare("insert ignore into seguidores (usuario_id, seguidor_id) values (?, ?)")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(userId, followerId); erro != nil {
		return erro
	}
	return nil

}

// Permite que um usuário pare de seguir outro usuário
func (repositories users) UnfollowUser(userId, followerId uint64) error {
	statement, erro := repositories.db.Prepare("delete from seguidores where usuario_id = ? and seguidor_id = ?")
	if erro != nil {
		return erro
	}

	if _, erro = statement.Exec(userId, followerId); erro != nil {
		return erro
	}
	return nil

}
