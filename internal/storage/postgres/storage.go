package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type storage struct {
	db  *pgx.Conn
	ctx context.Context
}

var (
	participantTableName = "participants_chat_"
	messageTableName     = "messages_chat_"
	userChatTableName    = "chats_user_"
)

// Connect подключается к базе данных и возвращает экземпляр хранилища
func Connect(ctx context.Context, connString string) (*storage, error) {
	db, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	s := &storage{
		db:  db,
		ctx: ctx,
	}

	return s, nil
}
