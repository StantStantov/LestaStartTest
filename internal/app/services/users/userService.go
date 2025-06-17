package users

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/services"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"context"
)

type UserService struct {
	userStore stores.UserStore

	idGen       services.IdGenerator
	psEncrypter services.PasswordEncrypter
	psValidator services.PasswordValidator
}

func NewUserService(
	userStore stores.UserStore,
	idGen services.IdGenerator,
	psEncrypter services.PasswordEncrypter,
	psValidator services.PasswordValidator,
) *UserService {
	return &UserService{
		userStore:   userStore,
		idGen:       idGen,
		psEncrypter: psEncrypter,
		psValidator: psValidator,
	}
}

func (s *UserService) IsValidUsername(username string) bool {
	if username == "" {
		return false
	}

	return true
}

func (s *UserService) IsValidPassword(password string) bool {
	if password == "" {
		return false
	}

	return true
}

func (s *UserService) Register(ctx context.Context, username, password string) error {
	user := s.newUser(username, password)

	if err := s.userStore.Register(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) IsRegistered(ctx context.Context, username, password string) (bool, error) {
	user, err := s.userStore.FindByName(ctx, username)
	if err != nil {
		return false, err
	}

	return user.IsUserPassword(password, s.psValidator), nil
}

func (s *UserService) newUser(username, password string) models.User {
	return models.NewUser(
		s.idGen.GenerateId(),
		username,
		s.psEncrypter.Encrypt(password),
	)
}
