package models_test

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"testing"
)

func TestCollectionAddDocument(t *testing.T) {
	t.Parallel()
	t.Run("PASS if added", func(t *testing.T) {
		t.Parallel()

		collection := models.NewEmptyCollection(0, "test", "123456")

		if err := collection.AddDocument(models.NewDocument(0, "file")); err != nil {
			t.Errorf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if duplicate", func(t *testing.T) {
		t.Parallel()

		collection := models.NewEmptyCollection(0, "test", "123456")

		term := models.NewDocument(0, "file")

		if err := collection.AddDocument(term); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := collection.AddDocument(term); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func TestCollectionFindDocument(t *testing.T) {
	t.Parallel()
	t.Run("PASS if found", func(t *testing.T) {
		t.Parallel()

		collection := models.NewEmptyCollection(0, "test", "123456")

		want := models.NewDocument(0, "file")

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

		collection := models.NewEmptyCollection(0, "test", "123456")

		if _, err := collection.FindDocument("Nothing"); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func TestCollectionRemoveDocument(t *testing.T) {
	t.Parallel()
	t.Run("PASS if removed", func(t *testing.T) {
		t.Parallel()

		collection := models.NewEmptyCollection(0, "test", "123456")

		want := models.NewDocument(0, "file")

		if err := collection.AddDocument(want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := collection.RemoveDocument(want.Name()); err != nil {
			t.Errorf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if empty", func(t *testing.T) {
		t.Parallel()

		collection := models.NewEmptyCollection(0, "test", "123456")

		if err := collection.RemoveDocument("Test"); err == nil {
			t.Errorf("Wanted %v, got %v", nil, err)
		}
	})
}
