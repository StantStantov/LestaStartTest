//go:build integration || !unit

package pgsql_test

import (
	"Stant/LestaGamesInternship/internal/infra/pgsql"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5"
)

func getTestTx(t *testing.T, db pgsql.DBConn, ctx context.Context) pgx.Tx {
	t.Helper()

	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			t.Fatal(err)
		}
	})

	return tx
}

func createTempDir(t *testing.T, dirpath, dirname string) string {
	t.Helper()

	dirpath = filepath.Join(dirpath, dirname)
	if err := os.RemoveAll(dirpath); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dirpath, 0770); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dirpath); err != nil {
		t.Fatal(err)
	}

	return dirpath
}

func createTempFile(t *testing.T, prefix string) *os.File {
	t.Helper()

	file, err := os.CreateTemp("", prefix+"_*")
	if err != nil {
		t.Fatalf("Wanted %v, got %v", nil, err)
	}

	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})

	return file
}
