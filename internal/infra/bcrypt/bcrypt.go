package bcrypt

import "golang.org/x/crypto/bcrypt"

func Encrypt(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash)
}

func ComparePasswords(hashedPassword, plainPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return false, err
	}
	return true, nil
}
