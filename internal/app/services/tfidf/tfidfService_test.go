//go:build unit || !integration

package tfidf_test

import (
	"Stant/LestaGamesInternship/internal/app/services/tfidf"
	"maps"
	"slices"
	"strings"
	"testing"
)

func TestProcessReaderToTerms(t *testing.T) {
	file := strings.NewReader("Random text for testing purposes only")

	want := []string{"Random", "text", "for", "testing", "purposes", "only"}
	got, err := tfidf.ProcessReaderToTerms(file)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, want) {
		t.Errorf("Wanted %+v, got %+v\n", want, got)
	}
}

func TestGetTermFrequency(t *testing.T) {
	text := []string{"word", "hello", "hello", "world"}

	want := map[string]uint64{"word": 1, "hello": 2, "world": 1}
	got := tfidf.GetTermFrequency(text)

	if !maps.Equal(got, want) {
		t.Errorf("Wanted %v, got %v\n", want, got)
	}
}

func TestCalculateIdf(t *testing.T) {
	termsAmount := uint64(10)
	termFrequency := uint64(1)

	want := 1.0
	got := tfidf.CalculateIdf(termsAmount, termFrequency)

	if want != got {
		t.Errorf("Wanted %f, got %f\n", want, got)
	}
}
