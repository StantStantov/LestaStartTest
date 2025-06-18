package services

type PasswordValidator interface {
	ComparePasswords(hashedPassword, plainPassword string) (bool, error)
}

type PasswordValidatorFunc func(string, string) (bool, error)

func (f PasswordValidatorFunc) ComparePasswords(hashedPassword, plainPassword string) (bool, error) {
	return f(hashedPassword, plainPassword)
}
