package postgres

import (
	"strconv"

	"github.com/MyLi2tlePony/messenger/internal/storage/entity"
)

func (s *storage) CreateParticipant(chatId int, p entity.Participant) (int, error) {
	query := `
		INSERT INTO
    		` + participantTableName + strconv.Itoa(chatId) + `
			(user_id, write, post, comment, delete, add_participant, delete_participant)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id`

	var id int
	err := s.db.QueryRow(s.ctx, query, p.UserId, p.Write, p.Post, p.Comment, p.Delete, p.AddParticipant, p.DeleteParticipant).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) SelectTopParticipants(chatId int, limit int) ([]entity.Participant, error) {
	query := `
		SELECT 
			id, user_id, write, post, comment, delete, add_participant, delete_participant
		FROM 
			` + participantTableName + strconv.Itoa(chatId) + `
		ORDER BY 
			id DESC
		LIMIT $1`

	rows, err := s.db.Query(s.ctx, query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	participants := make([]entity.Participant, 0, limit)
	for rows.Next() {
		var c entity.Participant

		if err = rows.Scan(&c.Id, &c.UserId, &c.Write, &c.Post, &c.Comment, &c.Delete, &c.AddParticipant, &c.DeleteParticipant); err != nil {
			return nil, err
		}

		participants = append(participants, c)
	}

	return participants, nil
}

func (s *storage) SelectParticipants(chatId int) ([]entity.Participant, error) {
	query := `
		SELECT 
			id, user_id, write, post, comment, delete, add_participant, delete_participant 
		FROM 
			` + participantTableName + strconv.Itoa(chatId)

	rows, err := s.db.Query(s.ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	participants := make([]entity.Participant, 0)
	for rows.Next() {
		var c entity.Participant

		if err = rows.Scan(&c.Id, &c.UserId, &c.Write, &c.Post, &c.Comment, &c.Delete, &c.AddParticipant, &c.DeleteParticipant); err != nil {
			return nil, err
		}

		participants = append(participants, c)
	}

	return participants, nil
}

func (s *storage) SelectParticipantIdByUserId(chatId int, userId int) (int, error) {
	query := `
		SELECT 
			id 
		FROM 
			` + participantTableName + strconv.Itoa(chatId) + `
		WHERE
			user_id = $1`

	row := s.db.QueryRow(s.ctx, query, userId)

	var id int
	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func (s *storage) DeleteParticipant(chatId int, id int) error {
	query := `DELETE FROM ` + participantTableName + strconv.Itoa(chatId) + ` WHERE id = $1`

	if _, err := s.db.Exec(s.ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (s *storage) CreateTableChatParticipants(chatId int) error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + participantTableName + strconv.Itoa(chatId) + `(
    		id SERIAL PRIMARY KEY UNIQUE NOT NULL,
			user_id INTEGER NOT NULL REFERENCES users (Id) ON DELETE CASCADE,
			
			write BOOLEAN NOT NULL,
			post BOOLEAN NOT NULL,
			comment BOOLEAN NOT NULL,
			delete BOOLEAN NOT NULL,
			add_participant BOOLEAN NOT NULL,
			delete_participant BOOLEAN NOT NULL
		)`

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) DeleteTableChatParticipants(chatId int) error {
	query := `DROP TABLE IF EXISTS ` + participantTableName + strconv.Itoa(chatId)

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}
