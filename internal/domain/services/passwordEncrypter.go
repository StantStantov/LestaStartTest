package services

type PasswordEncrypter interface {
	Encrypt(string) string
}

type PasswordEncrypterFunc func(string) string

func (f PasswordEncrypterFunc) Encrypt(plaintext string) string {
	return f(plaintext)
}
