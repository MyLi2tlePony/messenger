package postgres

import (
	"strconv"

	"github.com/MyLi2tlePony/messenger/internal/storage/entity"
)

// CreateUserChat добавляет идентификатор чата в чаты пользователя с указанным идентификатором
func (s *storage) CreateUserChat(userId, chatId int) (int, error) {
	query := `
		INSERT INTO
			` + userChatTableName + strconv.Itoa(userId) + `
			(chat_Id)
		VALUES
			($1)
		RETURNING 
			id`

	var id int

	err := s.db.QueryRow(s.ctx, query, chatId).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SelectUserChats возвращает все чаты пользователя по указанному идентификатору
func (s *storage) SelectUserChats(userId int) ([]entity.UserChats, error) {
	query := `
		SELECT 
			id, chat_id 
		FROM 
			` + userChatTableName + strconv.Itoa(userId)

	rows, err := s.db.Query(s.ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	couriers := make([]entity.UserChats, 0)
	for rows.Next() {
		var c entity.UserChats

		if err = rows.Scan(&c.Id, &c.ChatId); err != nil {
			return nil, err
		}

		couriers = append(couriers, c)
	}

	return couriers, nil
}

// DeleteUserChat удаляет чат с указанным идентификатором пользователя с указанным идентификатором
func (s *storage) DeleteUserChat(userId, userChatId int) error {
	query := `DELETE FROM ` + userChatTableName + strconv.Itoa(userId) + ` WHERE id = $1`

	if _, err := s.db.Exec(s.ctx, query, userChatId); err != nil {
		return err
	}

	return nil
}

// CreateTableUserChats создает таблицу с чатами для пользователя с указанным идентификатором
func (s *storage) CreateTableUserChats(userId int) error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + userChatTableName + strconv.Itoa(userId) + `(
    		id SERIAL PRIMARY KEY UNIQUE NOT NULL,
			chat_id INTEGER NOT NULL REFERENCES chats (Id) ON DELETE CASCADE
		)`

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTableUserChats удаляет таблицу с чатами пользователя с указанным идентификатором
func (s *storage) DeleteTableUserChats(userId int) error {
	query := `DROP TABLE IF EXISTS ` + userChatTableName + strconv.Itoa(userId)

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}
