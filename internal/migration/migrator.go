package migration

import (
	"context"
	"github.com/jackc/pgx/v4"
	"os"
	"time"
)

type migrator struct {
	db *pgx.Conn
}

func Connect(ctx context.Context, connString string) (*migrator, error) {
	db, err := pgx.Connect(ctx, connString)
	for i := 0; i < 5 && err != nil; i++ {
		time.Sleep(2 * time.Second)
		db, err = pgx.Connect(ctx, connString)
	}

	if err != nil {
		return nil, err
	}

	s := &migrator{
		db: db,
	}

	return s, nil
}

func (m *migrator) Migrate(ctx context.Context, path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(ctx, string(file))
	return err
}
