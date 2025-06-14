//go:build unit || !integration

package models_test

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/services"
	"crypto/rand"
	"testing"
)

func TestCollection(t *testing.T) {
	idGen := services.IdGeneratorFunc(func() string { return rand.Text() })

	t.Run("Test Add Document", func(t *testing.T) {
		t.Parallel()

		testCollectionAddDocument(t, idGen)
	})
	t.Run("Test Find Document", func(t *testing.T) {
		t.Parallel()

		testCollectionFindDocument(t, idGen)
	})
	t.Run("Test Remove Document", func(t *testing.T) {
		t.Parallel()

		testCollectionRemoveDocument(t, idGen)
	})
}

func testCollectionAddDocument(t *testing.T, idGen services.IdGenerator) {
	t.Helper()

	t.Run("PASS if added", func(t *testing.T) {
		t.Parallel()

		collection := models.NewEmptyCollection(idGen.GenerateId(), "test", "123456")

		if err := collection.AddDocument(models.NewDocument(idGen.GenerateId(), idGen.GenerateId(), "file", nil)); err != nil {
			t.Errorf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if duplicate", func(t *testing.T) {
		t.Parallel()

		collection := models.NewEmptyCollection(idGen.GenerateId(), "test", "123456")
		document := models.NewDocument(idGen.GenerateId(), idGen.GenerateId(), "file", nil)

		if err := collection.AddDocument(document); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := collection.AddDocument(document); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func testCollectionFindDocument(t *testing.T, idGen services.IdGenerator) {
	t.Helper()

	t.Run("PASS if found", func(t *testing.T) {
		t.Parallel()

		want := models.NewDocument(idGen.GenerateId(), idGen.GenerateId(), "file", nil)
		collection := models.NewEmptyCollection(idGen.GenerateId(), "test", "123456")

		if err := collection.AddDocument(want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		got, err := collection.FindDocument(want.Name())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}
	})
	t.Run("FAIL if not present", func(t *testing.T) {
		t.Parallel()

		collection := models.NewEmptyCollection(idGen.GenerateId(), "test", "123456")

		if _, err := collection.FindDocument("Nothing"); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func testCollectionRemoveDocument(t *testing.T, idGen services.IdGenerator) {
	t.Helper()

	t.Run("PASS if removed", func(t *testing.T) {
		t.Parallel()

		want := models.NewDocument(idGen.GenerateId(), idGen.GenerateId(), "file", nil)
		collection := models.NewEmptyCollection(idGen.GenerateId(), "test", "123456")

		if err := collection.AddDocument(want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := collection.RemoveDocument(want.Name()); err != nil {
			t.Errorf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if empty", func(t *testing.T) {
		t.Parallel()

		collection := models.NewEmptyCollection(idGen.GenerateId(), "test", "123456")

		if err := collection.RemoveDocument("Test"); err == nil {
			t.Errorf("Wanted %v, got %v", nil, err)
		}
	})
}
