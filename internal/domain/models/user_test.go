//go:build unit || !integration

package models_test

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/services"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestUser(t *testing.T) {
	idGen := services.IdGeneratorFunc(func() string { return rand.Text() })
	validator := services.PasswordValidatorFunc(mockPasswordValidator)
	encrypter := services.PasswordEncrypterFunc(func(s string) string { return s })

	t.Run("Test Check is User's password", func(t *testing.T) {
		t.Parallel()

		testUserIsUserPassword(t, idGen, validator, encrypter)
	})
	t.Run("Test Change Password", func(t *testing.T) {
		t.Parallel()

		testUserChangePassword(t, idGen, validator, encrypter)
	})
}

func testUserIsUserPassword(t *testing.T,
	idGen services.IdGenerator,
	validator services.PasswordValidator,
	encrypter services.PasswordEncrypter,
) {
	t.Helper()

	t.Run("PASS if same", func(t *testing.T) {
		t.Parallel()

		password := "password"
		user := models.NewUser(idGen.GenerateId(), "test", password, encrypter)

		if same := user.IsUserPassword(password, validator); !same {
			t.Errorf("Wanted %v, got %v", true, same)
		}
	})
	t.Run("FAIL if not same", func(t *testing.T) {
		t.Parallel()

		password := "password"
		anotherPassword := "yetAnotherPassword"
		user := models.NewUser(idGen.GenerateId(), "test", password, encrypter)

		if same := user.IsUserPassword(anotherPassword, validator); same {
			t.Errorf("Wanted %v, got %v", false, same)
		}
	})
}

func testUserChangePassword(t *testing.T,
	idGen services.IdGenerator,
	validator services.PasswordValidator,
	encrypter services.PasswordEncrypter,
) {
	t.Helper()

	t.Run("PASS if changed", func(t *testing.T) {
		t.Parallel()

		want := "newPassword"
		password := "password"
		user := models.NewUser(idGen.GenerateId(), "test", password, encrypter)

		if err := user.ChangePassword(password, want, validator, encrypter); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		got := user.HashedPassword()
		if want != got {
			t.Errorf("Wanted %s, got %s", want, got)
		}
	})
	t.Run("FAIL if incorrect current password", func(t *testing.T) {
		t.Parallel()

		password := "password"
		wrongPassword := "wrong"
		newPassword := "yetAnotherPassword"
		user := models.NewUser(idGen.GenerateId(), "test", password, encrypter)

		if err := user.ChangePassword(wrongPassword, newPassword, validator, encrypter); err == nil {
			t.Errorf("Wanted err, got %v", nil)
		}
	})
}

func mockPasswordValidator(s1, s2 string) (bool, error) {
	if s1 != s2 {
		return false, fmt.Errorf("%s != %s", s1, s2)
	}

	return true, nil
}
