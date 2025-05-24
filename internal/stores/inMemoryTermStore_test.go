package stores_test

import (
	"Stant/LestaGamesInternship/internal/models"
	"Stant/LestaGamesInternship/internal/stores"
	"slices"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	t.Run("Test Create", func(t *testing.T) {
		termStore := stores.NewInMemoryTermStore()
		testCreate(termStore, t)
	})
	t.Run("Test Read", func(t *testing.T) {
		termStore := stores.NewInMemoryTermStore()
		testRead(termStore, t)
	})
	t.Run("Test ReadAll", func(t *testing.T) {
		termStore := stores.NewInMemoryTermStore()
		testReadAll(termStore, t)
	})
	t.Run("Test CountAll", func(t *testing.T) {
		termStore := stores.NewInMemoryTermStore()
		testCountAll(termStore, t)
	})
	t.Run("Test Update", func(t *testing.T) {
		termStore := stores.NewInMemoryTermStore()
		testUpdate(termStore, t)
	})
	t.Run("Test Delete", func(t *testing.T) {
		termStore := stores.NewInMemoryTermStore()
		testDelete(termStore, t)
	})
}

func testCreate(termStore *stores.InMemoryTermStore, t *testing.T) {
	t.Helper()

	term := models.NewTerm("Test", 0, 0.0)

	if err := termStore.Create(term); err != nil {
		t.Fatal(err)
	}
}

func testRead(termStore *stores.InMemoryTermStore, t *testing.T) {
	t.Helper()

	want := models.NewTerm("Test", 0, 0.0)
	termStore.Create(want)

	got, err := termStore.Read(0)
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("Wanted %+v, got %+v", want, got)
	}
}

func testReadAll(termStore *stores.InMemoryTermStore, t *testing.T) {
	t.Helper()

	want := []models.Term{
		models.NewTerm("Romeo", 2, 1.57),
		models.NewTerm("salad", 2, 1.27),
		models.NewTerm("Falstaff", 4, 0.967),
	}
	for _, term := range want {
		termStore.Create(term)
	}

	got, err := termStore.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(want, got) {
		t.Errorf("Wanted %+v, got %+v", want, got)
	}
}

func testCountAll(termStore *stores.InMemoryTermStore, t *testing.T) {
	t.Helper()

	terms := []models.Term{
		models.NewTerm("Romeo", 2, 1.57),
		models.NewTerm("salad", 2, 1.27),
		models.NewTerm("Falstaff", 4, 0.967),
	}
	for _, term := range terms {
		termStore.Create(term)
	}

	want := 3
	got, err := termStore.CountAll()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("Wanted %+v, got %+v", want, got)
	}
}

func testUpdate(termStore *stores.InMemoryTermStore, t *testing.T) {
	t.Helper()

	term := models.NewTerm("Romeo", 2, 1.57)
	termStore.Create(term)
	newTerm := models.NewTerm("fool", 36, 0.012)
	termStore.Update(0, newTerm)

	want := newTerm
	got, err := termStore.Read(0)
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("Wanted %+v, got %+v", want, got)
	}
}

func testDelete(termStore *stores.InMemoryTermStore, t *testing.T) {
	t.Helper()

	terms := []models.Term{
		models.NewTerm("Romeo", 2, 1.57),
		models.NewTerm("salad", 2, 1.27),
		models.NewTerm("Falstaff", 4, 0.967),
	}
	for _, term := range terms {
		termStore.Create(term)
	}
	termStore.Delete(0)
	termStore.Delete(1)

	want := 1
	got, err := termStore.CountAll()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("Wanted %+v, got %+v", want, got)
	}
}
