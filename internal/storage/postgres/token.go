package postgres

// CreateToken создает токен с указанным текстом для пользователя с указанным идентификатором
func (s *storage) CreateToken(userId int, tokenText string) (int, error) {
	query := `
	INSERT INTO
		tokens (user_id, token)
	VALUES
		($1, $2)
	RETURNING 
		id`

	var id int
	err := s.db.QueryRow(s.ctx, query, userId, tokenText).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteToken удаляет токен с указанным идентификатором
func (s *storage) DeleteToken(id int) error {
	query := "DELETE FROM Tokens WHERE id = $1"

	if _, err := s.db.Exec(s.ctx, query, id); err != nil {
		return err
	}

	return nil
}
