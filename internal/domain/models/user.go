package models

import (
	"Stant/LestaGamesInternship/internal/domain/services"
	"fmt"
)

type User struct {
	id             string
	name           string
	hashedPassword string
}

func NewUser(id, name, password string) User {
	return User{
		id:             id,
		name:           name,
		hashedPassword: password,
	}
}

func (u *User) IsUserPassword(password string, validator services.PasswordValidator) bool {
	ok, err := validator.ComparePasswords(u.hashedPassword, password)
	if err != nil {
		return false
	}

	return ok
}

func (u *User) ChangePassword(currentPassword, newPassword string, validator services.PasswordValidator) error {
	ok, err := validator.ComparePasswords(u.hashedPassword, currentPassword)
	if err != nil {
		return fmt.Errorf("models/user.ChangePassword: [%w]", err)
	}
	if !ok {
		return fmt.Errorf("models/user.ChangePassword: [%v]", "Passwords are different")
	}
	u.hashedPassword = newPassword

	return nil
}

func (u *User) Id() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) HashedPassword() string {
	return u.hashedPassword
}
