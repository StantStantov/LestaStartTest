package models

import (
	"Stant/LestaGamesInternship/internal/domain/services"
	"fmt"
)

type User struct {
	uid            string
	name           string
	hashedPassword string
}

func NewUser(uid, name, password string) User {
	return User{
		uid:            uid,
		name:           name,
		hashedPassword: password,
	}
}

func (u *User) IsUserPassword(password string, validator services.PasswordValidator) bool {
	if err := validator.ComparePasswords(u.hashedPassword, password); err != nil {
		return false
	}

	return true
}

func (u *User) ChangePassword(currentPassword, newPassword string, validator services.PasswordValidator) error {
	if err := validator.ComparePasswords(u.hashedPassword, currentPassword); err != nil {
		return fmt.Errorf("models/user.ChangePassword: [%w]", err)
	}
	u.hashedPassword = newPassword

	return nil
}

func (u *User) Uid() string {
	return u.uid
}

func (u *User) Name() string {
	return u.name
}

func (u *User) HashedPassword() string {
	return u.hashedPassword
}
