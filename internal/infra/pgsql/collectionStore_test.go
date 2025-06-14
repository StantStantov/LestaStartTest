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
	"cmp"
	"context"
	"crypto/rand"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestCollectionStore(t *testing.T) {
	ctx := context.Background()
	idGen := services.IdGeneratorFunc(func() string { return rand.Text() })
	user := models.NewUser(idGen.GenerateId(), rand.Text(), rand.Text())
	file := apptest.CreateTestFile(t, "")
	document := models.NewDocument(idGen.GenerateId(), user.Id(), filepath.Base(file.Name()), file)

	dirpath := apptest.CreateTestDir(t, os.Getenv(config.PathToDocsEnv), "CollectionStore")
	fileStore := volume.NewFileStore(dirpath)

	dbPool := apptest.GetTestPool(t, ctx, os.Getenv(config.DatabaseUrlEnv))

	userStore := pgsql.NewUserStore(dbPool)
	if err := userStore.Register(ctx, user); err != nil {
		t.Fatal(err)
	}

	documentStore := pgsql.NewDocumentStore(dbPool, fileStore)
	if err := documentStore.Save(ctx, document); err != nil {
		t.Fatal(err)
	}

	t.Run("Test Save Collection", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)

		testCollectionStoreSave(t, ctx, tx, documentStore, idGen, user, document)
	})
	t.Run("Test Find Collection", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)

		testCollectionStoreFind(t, ctx, tx, documentStore, idGen, user, document)
	})
	t.Run("Test Pin Document to Collection", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)

		testCollectionStorePin(t, ctx, tx, documentStore, idGen, user, document)
	})
	t.Run("Test Unpin Document From Collection", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)

		testCollectionStoreUnpin(t, ctx, tx, documentStore, idGen, user, document)
	})
	t.Run("Test Rename Collection", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)

		testCollectionStoreRename(t, ctx, tx, documentStore, idGen, user, document)
	})
	t.Run("Test Delete Collection", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)

		testCollectionStoreDelete(t, ctx, tx, documentStore, idGen, user, document)
	})
}

func testCollectionStoreSave(
	t *testing.T,
	ctx context.Context,
	dbConn pgsql.DBConn,
	documentStore stores.DocumentStore,
	idGen services.IdGenerator,
	user models.User,
	document models.Document,
) {
	t.Helper()

	t.Run("PASS if saved", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		want := true
		id := idGen.GenerateId()
		userId := user.Id()
		name := rand.Text()
		documents := []models.Document{document}
		collection := models.NewCollection(id, userId, name, documents)

		if err := collectionStore.Save(ctx, *collection); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		got, err := collectionStore.IsExist(ctx, id)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}
	})
	t.Run("FAIL if duplicate", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		id := idGen.GenerateId()
		userId := user.Id()
		name := rand.Text()
		documents := []models.Document{document}
		collection := models.NewCollection(id, userId, name, documents)

		if err := collectionStore.Save(ctx, *collection); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := collectionStore.Save(ctx, *collection); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func testCollectionStoreFind(
	t *testing.T,
	ctx context.Context,
	dbConn pgsql.DBConn,
	documentStore stores.DocumentStore,
	idGen services.IdGenerator,
	user models.User,
	document models.Document,
) {
	t.Helper()

	isEqualDocuments := func(E1, E2 models.Document) bool {
		return (E1.Id() == E2.Id()) && (E1.UserId() == E2.UserId()) && (E1.Name() == E2.Name())
	}
	sortByCollectionId := func(E1, E2 *models.Collection) int {
		return cmp.Compare(E1.Id(), E2.Id())
	}

	t.Run("PASS if found by ID", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		wantId := idGen.GenerateId()
		wantUserId := user.Id()
		wantName := rand.Text()
		wantDocuments := []models.Document{document}
		want := models.NewCollection(wantId, wantUserId, wantName, wantDocuments)

		if err := collectionStore.Save(ctx, *want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		got, err := collectionStore.FindById(ctx, wantId)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		gotUserId := got.UserId()
		gotName := got.Name()
		gotDocuments := got.Documents()

		if wantUserId != gotUserId {
			t.Errorf("Wanted %v, got %v", wantUserId, gotUserId)
		}
		if wantName != gotName {
			t.Errorf("Wanted %v, got %v", wantName, gotName)
		}
		if !slices.EqualFunc(wantDocuments, gotDocuments, isEqualDocuments) {
			t.Errorf("Wanted %v, got %v", wantDocuments, gotDocuments)
		}
	})
	t.Run("PASS if found by User ID", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		wantUserId := user.Id()
		wantDocuments := []models.Document{document}
		want := []*models.Collection{
			models.NewCollection(idGen.GenerateId(), wantUserId, rand.Text(), wantDocuments),
			models.NewCollection(idGen.GenerateId(), wantUserId, rand.Text(), wantDocuments),
		}
		slices.SortFunc(want, sortByCollectionId)

		for _, collection := range want {
			if err := collectionStore.Save(ctx, *collection); err != nil {
				t.Fatalf("Wanted %v, got %v", nil, err)
			}
		}
		got, err := collectionStore.FindByUserId(ctx, wantUserId)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if len(want) != len(got) {
			t.Fatalf("Wanted len %v, got %v", len(want), len(got))
		}
		for i := range got {
			wantCollection := want[i]
			wantId := wantCollection.Id()
			wantName := wantCollection.Name()
			wantDocuments := wantCollection.Documents()

			gotCollection := got[i]
			gotId := gotCollection.Id()
			gotName := gotCollection.Name()
			gotDocuments := gotCollection.Documents()

			if wantId != gotId {
				t.Errorf("Wanted ID %v, got %v", wantId, gotId)
			}
			if wantName != gotName {
				t.Errorf("Wanted Name %v, got %v", wantName, gotName)
			}
			if !slices.EqualFunc(wantDocuments, gotDocuments, isEqualDocuments) {
				t.Errorf("Wanted Documents %v, got %v", wantDocuments, gotDocuments)
			}
		}
	})
	t.Run("FAIL if ID doesn't exist", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		if _, err := collectionStore.FindById(ctx, ""); err == nil {
			t.Fatalf("Wanted err, got %v", err)
		}
	})
	t.Run("FAIL if User ID doesn't have collections", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		if _, err := collectionStore.FindByUserId(ctx, ""); err == nil {
			t.Fatalf("Wanted err, got %v", err)
		}
	})
}

func testCollectionStorePin(
	t *testing.T,
	ctx context.Context,
	dbConn pgsql.DBConn,
	documentStore stores.DocumentStore,
	idGen services.IdGenerator,
	user models.User,
	document models.Document,
) {
	t.Helper()

	t.Run("PASS if pinned", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		want := true
		id := idGen.GenerateId()
		userId := user.Id()
		name := rand.Text()
		documents := []models.Document{}
		collection := models.NewCollection(id, userId, name, documents)

		if err := collectionStore.Save(ctx, *collection); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := collectionStore.PinDocument(ctx, collection.Id(), document.Id()); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		got, err := collectionStore.IsPinned(ctx, collection.Id(), document.Id())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}
	})
	t.Run("FAIL if Document already pinned", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		id := idGen.GenerateId()
		userId := user.Id()
		name := rand.Text()
		documents := []models.Document{document}
		collection := models.NewCollection(id, userId, name, documents)

		if err := collectionStore.Save(ctx, *collection); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := collectionStore.PinDocument(ctx, collection.Id(), document.Id()); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
	t.Run("FAIL if Collection doesn't exist", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		if err := collectionStore.PinDocument(ctx, "", document.Id()); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func testCollectionStoreUnpin(
	t *testing.T,
	ctx context.Context,
	dbConn pgsql.DBConn,
	documentStore stores.DocumentStore,
	idGen services.IdGenerator,
	user models.User,
	document models.Document,
) {
	t.Helper()

	t.Run("PASS if unpinned", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		want := false
		id := idGen.GenerateId()
		userId := user.Id()
		name := rand.Text()
		documents := []models.Document{document}
		collection := models.NewCollection(id, userId, name, documents)

		if err := collectionStore.Save(ctx, *collection); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := collectionStore.UnpinDocument(ctx, collection.Id(), document.Id()); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		got, err := collectionStore.IsPinned(ctx, collection.Id(), document.Id())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}
	})
	t.Run("FAIL if Collection doesn't exist", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		if err := collectionStore.UnpinDocument(ctx, "", document.Id()); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
	t.Run("FAIL if Document doesn't exist", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		id := idGen.GenerateId()
		userId := user.Id()
		name := rand.Text()
		documents := []models.Document{document}
		collection := models.NewCollection(id, userId, name, documents)

		if err := collectionStore.Save(ctx, *collection); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := collectionStore.UnpinDocument(ctx, collection.Id(), ""); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func testCollectionStoreRename(
	t *testing.T,
	ctx context.Context,
	dbConn pgsql.DBConn,
	documentStore stores.DocumentStore,
	idGen services.IdGenerator,
	user models.User,
	document models.Document,
) {
	t.Helper()

	t.Run("PASS if renamed", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		id := idGen.GenerateId()
		userId := user.Id()
		name := rand.Text()
		wantName := "Renamed" + name
		documents := []models.Document{document}
		collection := models.NewCollection(id, userId, name, documents)

		if err := collectionStore.Save(ctx, *collection); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := collectionStore.Rename(ctx, collection.Id(), wantName); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if doesn't exist", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		if err := collectionStore.Rename(ctx, "", rand.Text()); err == nil {
			t.Fatalf("Wanted err, got %v", err)
		}
	})
}

func testCollectionStoreDelete(
	t *testing.T,
	ctx context.Context,
	dbConn pgsql.DBConn,
	documentStore stores.DocumentStore,
	idGen services.IdGenerator,
	user models.User,
	document models.Document,
) {
	t.Helper()

	t.Run("PASS if deleted", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		want := false
		id := idGen.GenerateId()
		userId := user.Id()
		name := rand.Text()
		documents := []models.Document{document}
		collection := models.NewCollection(id, userId, name, documents)

		if err := collectionStore.Save(ctx, *collection); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := collectionStore.Delete(ctx, id); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		gotIsExist, err := collectionStore.IsExist(ctx, id)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if want != gotIsExist {
			t.Errorf("Wanted %v, got %v", want, gotIsExist)
		}
		gotIsPinned, err := collectionStore.IsPinned(ctx, collection.Id(), document.Id())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if want != gotIsPinned {
			t.Errorf("Wanted %v, got %v", want, gotIsPinned)
		}
	})
	t.Run("FAIL if doesn't exist", func(t *testing.T) {
		tx := apptest.GetTestTx(t, ctx, dbConn)
		collectionStore := pgsql.NewCollectionStore(tx, documentStore)

		if err := collectionStore.Delete(ctx, ""); err == nil {
			t.Fatalf("Wanted err, got %v", err)
		}
	})
}
