package models_test

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"crypto/rand"
	"testing"
)

func TestStatisticsAddTerm(t *testing.T) {
	t.Parallel()
	t.Run("PASS if added", func(t *testing.T) {
		t.Parallel()

		stat := models.NewEmptyStatistic(0)

		want := models.NewTerm("Test", 1, 1)

		if err := stat.AddTerm(want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		got, err := stat.FindTerm(want.Word())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}
	})
	t.Run("FAIL if max capacity", func(t *testing.T) {
		t.Parallel()
		stat := models.NewEmptyStatistic(0)

		for range models.MaxStatisticTermsAmount {
			if err := stat.AddTerm(models.NewTerm(rand.Text(), 1, 1)); err != nil {
				t.Fatalf("Wanted %v, got %v", nil, err)
			}
		}

		if err := stat.AddTerm(models.NewTerm(rand.Text(), 1, 1)); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
	t.Run("FAIL if duplicate", func(t *testing.T) {
		t.Parallel()

		stat := models.NewEmptyStatistic(0)

		term := models.NewTerm("Test", 1, 1)

		if err := stat.AddTerm(term); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if err := stat.AddTerm(term); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func TestStatisticsFindTerm(t *testing.T) {
	t.Parallel()
	t.Run("PASS if found", func(t *testing.T) {
		t.Parallel()

		stat := models.NewEmptyStatistic(0)

		want := models.NewTerm("Test", 1, 1)

		if err := stat.AddTerm(want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		got, err := stat.FindTerm(want.Word())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}
	})
	t.Run("FAIL if not present", func(t *testing.T) {
		t.Parallel()

		stat := models.NewEmptyStatistic(0)

		_, err := stat.FindTerm("Test")
		if err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func TestStatisticsRemoveTerm(t *testing.T) {
	t.Parallel()
	t.Run("PASS if removed", func(t *testing.T) {
		t.Parallel()

		stat := models.NewEmptyStatistic(0)

		term := models.NewTerm("Test", 1, 1)

		if err := stat.AddTerm(term); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		err := stat.RemoveTerm(term.Word())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		_, err = stat.FindTerm(term.Word())
		if err != nil {
			t.Errorf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if empty", func(t *testing.T) {
		t.Parallel()

		stat := models.NewEmptyStatistic(0)

		if err := stat.RemoveTerm("Test"); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}
