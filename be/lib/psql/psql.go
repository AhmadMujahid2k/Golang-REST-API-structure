package psql

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tx struct {
	Tx pgx.Tx
}
type Postgres struct {
	Db *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
	err        error
)

func Init(ctx context.Context, connString string) (*Postgres, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			err = fmt.Errorf("unable to create connection pool: %w", err)
			return
		}

		pgInstance = &Postgres{db}
	})

	return pgInstance, err
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.Db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.Db.Close()
}

func QueryRow[T any](pg *Postgres, tx *Tx, query string, args ...any) (*T, error) {
	var ret T
	var rows pgx.Rows
	var err error

	ctx := context.Background()

	if tx != nil {
		rows, err = tx.Tx.Query(ctx, query, args...)
	} else {
		rows, err = pg.Db.Query(ctx, query, args...)
	}

	if err != nil {
		return &ret, err
	}
	defer rows.Close()

	ret, err = pgx.CollectOneRow[T](rows, pgx.RowToStructByName[T])
	if err != nil {
		return &ret, err
	}
	return &ret, nil
}

func Query[T any](pg *Postgres, tx *Tx, query string, args ...any) (*[]T, error) {
	var rows pgx.Rows
	var err error

	ctx := context.Background()

	if tx != nil {
		rows, err = tx.Tx.Query(ctx, query, args...)
	} else {
		rows, err = pg.Db.Query(ctx, query, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dest, err := pgx.CollectRows[T](rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, err
	}
	return &dest, nil
}

func Exec(pg *Postgres, tx *Tx, query string, args ...any) error {
	ctx := context.Background()
	var err error

	if tx != nil {
		_, err = tx.Tx.Exec(ctx, query, args...)
	} else {
		_, err = pg.Db.Exec(ctx, query, args...)
	}

	return err
}

func (pg *Postgres) Begin() (*Tx, error) {
	tx, err := pg.Db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx}, nil
}

func (tx *Tx) Rollback() error {
	return tx.Tx.Rollback(context.Background())
}

func (tx *Tx) Commit() error {
	return tx.Tx.Commit(context.Background())
}

