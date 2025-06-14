package services

type PasswordValidator interface {
	ComparePasswords(hashedPassword, plainPassword string) error
}

type PasswordValidatorFunc func(string, string) error

func (f PasswordValidatorFunc) ComparePasswords(hashedPassword, plainPassword string) error {
	return f(hashedPassword, plainPassword)
}
