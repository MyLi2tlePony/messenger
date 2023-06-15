package migration

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"time"
)

type migrator struct {
	db *pgx.Conn
}

// Connect пытается подключиться к базе данных 5 раз с задержкой в 2 секунды
// это необходимо если приложение в контейнере запустилось быстрее базы данных
func Connect(ctx context.Context, connString string) (*migrator, error) {
	fmt.Println("try db connect")
	db, err := pgx.Connect(ctx, connString)

	for i := 0; i < 5 && err != nil; i++ {
		time.Sleep(2 * time.Second)

		fmt.Println("try db connect")
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

// Migrate применение миграции из файла
func (m *migrator) Migrate(ctx context.Context, path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(ctx, string(file))
	return err
}
