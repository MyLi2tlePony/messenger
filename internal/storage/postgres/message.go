package postgres

import (
	"strconv"

	"github.com/MyLi2tlePony/messenger/internal/storage/entity"
)

func (s *storage) CreateMessage(chatId int, chat entity.Message) (int, error) {
	query := `
		INSERT INTO
    		` + messageTableName + strconv.Itoa(chatId) + `
			(user_id, commented_message_id, message, changed, created)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			id`

	var id int
	err := s.db.
		QueryRow(s.ctx, query, chat.UserId, chat.CommentedMessageId, chat.Text, chat.Changed, chat.Created).
		Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) SelectTopMessages(chatId, limit int) ([]entity.Message, error) {
	query := `
		SELECT 
			id, user_id, commented_message_id, message, changed, created 
		FROM 
			` + messageTableName + strconv.Itoa(chatId) + `
		ORDER BY 
			id DESC
		LIMIT $1`

	rows, err := s.db.Query(s.ctx, query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	couriers := make([]entity.Message, 0, limit)
	for rows.Next() {
		var c entity.Message

		if err = rows.Scan(&c.Id, &c.UserId, &c.CommentedMessageId, &c.Text, &c.Changed, &c.Created); err != nil {
			return nil, err
		}

		couriers = append(couriers, c)
	}

	return couriers, nil
}

func (s *storage) SelectMessagesById(chatId, minId, maxId int) ([]entity.Message, error) {
	query := `
		SELECT 
			id, user_id, commented_message_id, message, changed, created 
		FROM 
			` + messageTableName + strconv.Itoa(chatId) + `
		WHERE
			id >= $1 AND id <= $2`

	rows, err := s.db.Query(s.ctx, query, minId, maxId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	couriers := make([]entity.Message, 0, maxId-minId+1)
	for rows.Next() {
		var c entity.Message

		if err = rows.Scan(&c.Id, &c.UserId, &c.CommentedMessageId, &c.Text, &c.Changed, &c.Created); err != nil {
			return nil, err
		}

		couriers = append(couriers, c)
	}

	return couriers, nil
}

func (s *storage) DeleteMessage(chatId int, id int) error {
	query := `DELETE FROM ` + messageTableName + strconv.Itoa(chatId) + ` WHERE id = $1`

	if _, err := s.db.Exec(s.ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (s *storage) CreateTableChatMessages(chatId int) error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + messageTableName + strconv.Itoa(chatId) + `(
    		id SERIAL PRIMARY KEY UNIQUE NOT NULL,
			user_id INTEGER NOT NULL REFERENCES users (Id) ON DELETE CASCADE,
			
			commented_message_id INTEGER,
			message TEXT NOT NULL,
			changed BOOLEAN NOT NULL,

			created TIMESTAMP NOT NULL
		)`

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) DeleteTableChatMessages(chatId int) error {
	query := `DROP TABLE IF EXISTS ` + messageTableName + strconv.Itoa(chatId)

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}
