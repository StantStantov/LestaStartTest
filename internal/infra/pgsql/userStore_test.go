//go:build integration || !unit

package pgsql_test

import (
	"Stant/LestaGamesInternship/internal/app/config"
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/services"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"Stant/LestaGamesInternship/internal/infra/pgsql"
	"Stant/LestaGamesInternship/internal/pkg/apptest"
	"context"
	"crypto/rand"
	"os"
	"testing"
)

func TestUserStore(t *testing.T) {
	ctx := context.Background()

	dbPool := apptest.GetTestPool(t, ctx, os.Getenv(config.DatabaseUrlEnv))
	encrypter := services.PasswordEncrypterFunc(func(s string) string { return s })

	t.Run("Test Register User", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)
		userStore := pgsql.NewUserStore(tx, encrypter)

		testUserStoreRegister(t, ctx, userStore, encrypter)
	})
	t.Run("Test Find User", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)
		userStore := pgsql.NewUserStore(tx, encrypter)

		testUserStoreFind(t, ctx, userStore, encrypter)
	})
	t.Run("Test Update User", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)
		userStore := pgsql.NewUserStore(tx, encrypter)

		testUserStoreUpdate(t, ctx, userStore, encrypter)
	})
	t.Run("Test Deregister User", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)
		userStore := pgsql.NewUserStore(tx, encrypter)

		testUserStoreDeregister(t, ctx, userStore, encrypter)
	})
}

func testUserStoreRegister(t *testing.T, ctx context.Context, userStore stores.UserStore, encrypter services.PasswordEncrypter) {
	t.Helper()

	t.Run("PASS if registered", func(t *testing.T) {
		want := true
		user := models.NewUser(rand.Text(), rand.Text(), rand.Text(), encrypter)

		if err := userStore.Register(ctx, user); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		gotById, err := userStore.IsIdRegistered(ctx, user.Id())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		gotByName, err := userStore.IsNameRegistered(ctx, user.Name())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if want != gotById {
			t.Errorf("Wanted %v, got %v", want, gotById)
		}
		if want != gotByName {
			t.Errorf("Wanted %v, got %v", want, gotByName)
		}
	})
	t.Run("FAIL if already present", func(t *testing.T) {
		user := models.NewUser(rand.Text(), rand.Text(), rand.Text(), encrypter)

		if err := userStore.Register(ctx, user); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := userStore.Register(ctx, user); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func testUserStoreFind(t *testing.T, ctx context.Context, userStore stores.UserStore, encrypter services.PasswordEncrypter) {
	t.Helper()

	t.Run("PASS if found", func(t *testing.T) {
		want := models.NewUser(rand.Text(), rand.Text(), rand.Text(), encrypter)

		if err := userStore.Register(ctx, want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		got, err := userStore.FindById(ctx, want.Id())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}
	})
	t.Run("FAIL if not present", func(t *testing.T) {
		want := models.User{}

		got, err := userStore.FindById(ctx, rand.Text())
		if err == nil {
			t.Errorf("Wanted err, got %v", err)
		}

		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}
	})
}

func testUserStoreUpdate(t *testing.T, ctx context.Context, userStore stores.UserStore, encrypter services.PasswordEncrypter) {
	t.Helper()

	mockValidator := services.PasswordValidatorFunc(func(s1, s2 string) (bool, error) { return true, nil })

	t.Run("PASS if updated", func(t *testing.T) {
		want := models.NewUser(rand.Text(), rand.Text(), rand.Text(), encrypter)

		if err := userStore.Register(ctx, want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		want.ChangePassword(want.HashedPassword(), rand.Text(), mockValidator, encrypter)
		if err := userStore.Update(ctx, want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if not present", func(t *testing.T) {
		want := models.NewUser(rand.Text(), rand.Text(), rand.Text(), encrypter)

		if err := userStore.Update(ctx, want); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func testUserStoreDeregister(t *testing.T, ctx context.Context, userStore stores.UserStore, encrypter services.PasswordEncrypter) {
	t.Helper()

	t.Run("PASS if deregistered", func(t *testing.T) {
		want := false
		user := models.NewUser(rand.Text(), rand.Text(), rand.Text(), encrypter)

		if err := userStore.Register(ctx, user); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := userStore.Deregister(ctx, user.Id()); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		gotById, err := userStore.IsIdRegistered(ctx, user.Id())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		gotByName, err := userStore.IsNameRegistered(ctx, user.Name())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if want != gotById {
			t.Errorf("Wanted %v, got %v", want, gotById)
		}
		if want != gotByName {
			t.Errorf("Wanted %v, got %v", want, gotByName)
		}
	})
	t.Run("FAIL if not present", func(t *testing.T) {
		if err := userStore.Deregister(ctx, rand.Text()); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}
