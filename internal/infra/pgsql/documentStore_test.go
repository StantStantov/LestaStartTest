//go:build integration || !unit

package pgsql_test

import (
	"Stant/LestaGamesInternship/internal/app/config"
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/services"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"Stant/LestaGamesInternship/internal/infra/pgsql"
	"Stant/LestaGamesInternship/internal/infra/volume"
	"Stant/LestaGamesInternship/internal/pkg/apptest"
	"context"
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocumentStore(t *testing.T) {
	ctx := context.Background()

	idGen := services.IdGeneratorFunc(func() string { return rand.Text() })
	encrypter := services.PasswordEncrypterFunc(func(s string) string { return s })
	user := models.NewUser(idGen.GenerateId(), rand.Text(), rand.Text(), encrypter)
	dbPool := apptest.GetTestPool(t, ctx, os.Getenv(config.DatabaseUrlEnv))

	dirpath := apptest.CreateTestDir(t, os.Getenv(config.PathToDocsEnv), "DocumentStore")
	fileStore := volume.NewFileStore(dirpath)

	t.Run("Test Save Document", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)

		userStore := pgsql.NewUserStore(tx, encrypter)
		if err := userStore.Register(ctx, user); err != nil {
			t.Fatal(err)
		}

		documentStore := pgsql.NewDocumentStore(tx, fileStore)

		testDocumentStoreSave(t, ctx, documentStore, user, idGen)
	})
	t.Run("Test Open Document", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)

		userStore := pgsql.NewUserStore(tx, encrypter)
		if err := userStore.Register(ctx, user); err != nil {
			t.Fatal(err)
		}

		documentStore := pgsql.NewDocumentStore(tx, fileStore)

		testDocumentStoreOpen(t, ctx, documentStore, user, idGen)
	})
	t.Run("Test Rename Document", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)

		userStore := pgsql.NewUserStore(tx, encrypter)
		if err := userStore.Register(ctx, user); err != nil {
			t.Fatal(err)
		}

		documentStore := pgsql.NewDocumentStore(tx, fileStore)

		testDocumentStoreRename(t, ctx, documentStore, user, idGen)
	})
	t.Run("Test Delete Document", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)

		userStore := pgsql.NewUserStore(tx, encrypter)
		if err := userStore.Register(ctx, user); err != nil {
			t.Fatal(err)
		}

		documentStore := pgsql.NewDocumentStore(tx, fileStore)

		testDocumentStoreDelete(t, ctx, documentStore, user, idGen)
	})
}

func testDocumentStoreSave(t *testing.T,
	ctx context.Context,
	documentStore stores.DocumentStore,
	user models.User,
	idGen services.IdGenerator,
) {
	t.Helper()

	t.Run("PASS if saved", func(t *testing.T) {
		want := true
		file := apptest.CreateTestFile(t, "")
		file.WriteString("hello, documents")
		file.Seek(0, 0)
		id := idGen.GenerateId()
		userId := user.Id()
		filename := filepath.Base(file.Name())
		document := models.NewDocument(id, userId, filename, file)

		if err := documentStore.Save(ctx, document); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		gotById, err := documentStore.IsIdExist(ctx, id)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		gotByName, err := documentStore.IsNameExist(ctx, userId, filename)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if want != gotById {
			t.Fatalf("Wanted %v, got %v", want, gotById)
		}
		if want != gotByName {
			t.Fatalf("Wanted %v, got %v", want, gotByName)
		}
	})
	t.Run("FAIL if already exists", func(t *testing.T) {
		file := apptest.CreateTestFile(t, "")
		id := idGen.GenerateId()
		userId := user.Id()
		filename := filepath.Base(file.Name())
		document := models.NewDocument(id, userId, filename, file)

		if err := documentStore.Save(ctx, document); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := documentStore.Save(ctx, document); err == nil {
			t.Fatalf("Wanted err, got %v", err)
		}
	})
}

func testDocumentStoreOpen(t *testing.T,
	ctx context.Context,
	documentStore stores.DocumentStore,
	user models.User,
	idGen services.IdGenerator,
) {
	t.Helper()

	t.Run("PASS if opened", func(t *testing.T) {
		strBuilder := new(strings.Builder)
		file := apptest.CreateTestFile(t, "")
		wantId := idGen.GenerateId()
		wantUserId := user.Id()
		wantName := filepath.Base(file.Name())
		wantData := ""
		document := models.NewDocument(wantId, wantUserId, wantName, file)

		if err := documentStore.Save(ctx, document); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		got, err := documentStore.Open(ctx, wantId)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		io.Copy(strBuilder, got.File())
		gotId := got.Id()
		gotName := got.Name()
		gotData := strBuilder.String()

		if wantId != gotId {
			t.Errorf("Wanted %v, got %v", wantId, gotId)
		}
		if wantName != gotName {
			t.Errorf("Wanted %v, got %v", wantName, gotName)
		}
		if wantData != gotData {
			t.Errorf("Wanted %v, got %v", wantData, gotData)
		}
	})
	t.Run("FAIL if doesn't exist", func(t *testing.T) {
		if _, err := documentStore.Open(ctx, idGen.GenerateId()); err == nil {
			t.Fatalf("Wanted err, got %v", err)
		}
	})
}

func testDocumentStoreRename(t *testing.T,
	ctx context.Context,
	documentStore stores.DocumentStore,
	user models.User,
	idGen services.IdGenerator,
) {
	t.Helper()

	t.Run("PASS if renamed", func(t *testing.T) {
		file := apptest.CreateTestFile(t, "")
		id := idGen.GenerateId()
		userId := user.Id()
		fileName := filepath.Base(file.Name())
		wantName := "Renamed" + fileName
		document := models.NewDocument(id, userId, fileName, file)

		if err := documentStore.Save(ctx, document); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := documentStore.Rename(ctx, id, wantName); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if doesn't exists", func(t *testing.T) {
		if err := documentStore.Rename(ctx, rand.Text(), rand.Text()); err == nil {
			t.Fatalf("Wanted err, got %v", err)
		}
	})
}

func testDocumentStoreDelete(t *testing.T,
	ctx context.Context,
	documentStore stores.DocumentStore,
	user models.User,
	idGen services.IdGenerator,
) {
	t.Helper()

	t.Run("PASS if deleted", func(t *testing.T) {
		file := apptest.CreateTestFile(t, "")
		id := idGen.GenerateId()
		userId := user.Id()
		fileName := filepath.Base(file.Name())
		document := models.NewDocument(id, userId, fileName, file)

		if err := documentStore.Save(ctx, document); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := documentStore.Delete(ctx, id); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if doesn't exist", func(t *testing.T) {
		if err := documentStore.Delete(ctx, rand.Text()); err == nil {
			t.Fatalf("Wanted err, got %v", err)
		}
	})
}
