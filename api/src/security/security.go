package security

import "golang.org/x/crypto/bcrypt"

// Recebe a string e aplica o hash
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

}

// Faz a verificação da senha
func Verify(passwordString, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordString))

}
