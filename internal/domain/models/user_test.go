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

	t.Run("Test Check is User's password", func(t *testing.T) {
		t.Parallel()

		testUserIsUserPassword(t, idGen, validator)
	})
	t.Run("Test Change Password", func(t *testing.T) {
		t.Parallel()

		testUserChangePassword(t, idGen, validator)
	})
}

func testUserIsUserPassword(t *testing.T,
	idGen services.IdGenerator,
	validator services.PasswordValidator,
) {
	t.Helper()

	t.Run("PASS if same", func(t *testing.T) {
		t.Parallel()

		password := "password"
		user := models.NewUser(idGen.GenerateId(), "test", password)

		if same := user.IsUserPassword(password, validator); !same {
			t.Errorf("Wanted %v, got %v", true, same)
		}
	})
	t.Run("FAIL if not same", func(t *testing.T) {
		t.Parallel()

		password := "password"
		anotherPassword := "yetAnotherPassword"
		user := models.NewUser(idGen.GenerateId(), "test", password)

		if same := user.IsUserPassword(anotherPassword, validator); same {
			t.Errorf("Wanted %v, got %v", false, same)
		}
	})
}

func testUserChangePassword(t *testing.T,
	idGen services.IdGenerator,
	validator services.PasswordValidator,
) {
	t.Helper()

	t.Run("PASS if changed", func(t *testing.T) {
		t.Parallel()

		want := "newPassword"
		password := "password"
		user := models.NewUser(idGen.GenerateId(), "test", password)

		if err := user.ChangePassword(password, want, validator); err != nil {
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
		user := models.NewUser(idGen.GenerateId(), "test", password)

		if err := user.ChangePassword(wrongPassword, newPassword, validator); err == nil {
			t.Errorf("Wanted err, got %v", nil)
		}
	})
}

func mockPasswordValidator(s1, s2 string) error {
	if s1 != s2 {
		return fmt.Errorf("%s != %s", s1, s2)
	}

	return nil
}
