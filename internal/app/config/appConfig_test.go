//go:build unit || !integration

package config

import (
	"os"
	"testing"
)

func TestReadServerPort(t *testing.T) {
	t.Parallel()

	t.Run("PASS if port is correct", func(t *testing.T) {
		t.Parallel()

		want := "9090"

		err := os.Setenv(ServerPortEnv, want)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		got, err := readServerPort()
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if want != got {
			t.Errorf("Wanted %s, got %s", want, got)
		}
	})
	t.Run("FAIL if port not set", func(t *testing.T) {
		t.Parallel()

		if err := os.Unsetenv(ServerPortEnv); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		_, err := readServerPort()
		if err == nil {
			t.Fatalf("Wanted err, got nil")
		}
	})
	t.Run("FAIL if port is incorrect", func(t *testing.T) {
		t.Parallel()

		ports := []string{"", "2000000", "tseta"}

		for _, port := range ports {
			err := os.Setenv(ServerPortEnv, port)
			if err != nil {
				t.Fatalf("Wanted %v, got %v", nil, err)
			}

			_, err = readServerPort()
			if err == nil {
				t.Errorf("Wanted err, got %v", err)
			}
		}
	})
}

func TestReadDbUrl(t *testing.T) {
	t.Parallel()

	t.Run("PASS if read", func(t *testing.T) {
		t.Parallel()

		want := "postgres://john:secret@localhost:5432/db"

		if err := os.Setenv(DatabaseUrlEnv, want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		got, err := readDatabaseUrl()
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if want != got {
			t.Errorf("Wanted %s, got %s", want, got)
		}
	})
	t.Run("FAIL if not set", func(t *testing.T) {
		t.Parallel()

		if err := os.Unsetenv(DatabaseUrlEnv); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		_, err := readDatabaseUrl()
		if err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func TestReadPathToDocs(t *testing.T) {
	t.Parallel()

	t.Run("PASS if read", func(t *testing.T) {
		t.Parallel()

		want := "/documents"

		if err := os.Setenv(PathToDocsEnv, want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		got, err := readPathToDocuments()
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if want != got {
			t.Errorf("Wanted %s, got %s", want, got)
		}
	})
	t.Run("FAIL if not set", func(t *testing.T) {
		t.Parallel()

		if err := os.Unsetenv(PathToDocsEnv); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		_, err := readPathToDocuments()
		if err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}
