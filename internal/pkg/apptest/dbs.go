package apptest

import (
	"Stant/LestaGamesInternship/internal/infra/pgsql"
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetTestPool(t *testing.T, ctx context.Context, dbUrl string) *pgxpool.Pool {
	t.Helper()

	dbPool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		dbPool.Close()
	})

	return dbPool
}

func GetTestTx(t *testing.T, ctx context.Context, db pgsql.DBConn) pgx.Tx {
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
