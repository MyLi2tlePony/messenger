package postgres

import "github.com/MyLi2tlePony/messenger/internal/storage/entity"

func (s *storage) CreateChat(chat entity.Chat) (int, error) {
	query := `
		INSERT INTO
    		chats (name, description, open, created)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			id`

	var id int
	err := s.db.QueryRow(s.ctx, query, chat.Name, chat.Description, chat.Open, chat.Created).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) SelectChatById(id int) (entity.Chat, error) {
	query := `
		SELECT 
			id, name, description, open, created
		FROM 
			chats 
		WHERE 
			id = $1`

	row := s.db.QueryRow(s.ctx, query, id)
	chat := entity.Chat{}

	if err := row.Scan(&chat.Id, &chat.Name, &chat.Description, &chat.Open, &chat.Created); err != nil {
		return entity.Chat{}, err
	}

	return chat, nil
}

func (s *storage) DeleteChat(id int) error {
	query := "DELETE FROM chats WHERE id = $1"

	if _, err := s.db.Exec(s.ctx, query, id); err != nil {
		return err
	}

	return nil
}
