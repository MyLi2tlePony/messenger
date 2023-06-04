package postgres

import (
	"context"
	"github.com/MyLi2tlePony/messenger/internal/storage/entity"
	"strconv"

	"github.com/jackc/pgx/v4"
)

type storage struct {
	db  *pgx.Conn
	ctx context.Context
}

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

func (s *storage) CreateUser(login, password string) (int, error) {
	query := `
		INSERT INTO
    		users (login, password, first_name, second_name)
		VALUES
    		($1, $2, '', '')
		RETURNING
			id`

	var id int
	err := s.db.QueryRow(s.ctx, query, login, password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = $1"

	if _, err := s.db.Exec(s.ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (s *storage) SelectUserByPublicId(publicId string) (entity.User, error) {
	query := "SELECT id, login, password, first_name, second_name FROM users WHERE public_id = $1"

	row := s.db.QueryRow(s.ctx, query, publicId)
	user := entity.User{
		PublicId: publicId,
	}

	if err := row.Scan(&user.Id, &user.Login, &user.Password, &user.FirstName, &user.SecondName); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *storage) SelectUserByTocken(tocken entity.Tocken) (entity.User, error) {
	query := `
		SELECT 
			id, public_id, login, password, first_name, second_name 
		FROM 
			users 
		WHERE 
			id = (SELECT user_Id FROM tockens WHERE tocken = $1)`

	rows := s.db.QueryRow(s.ctx, query, tocken.Text)

	var user entity.User
	if err := rows.Scan(&user.Id, &user.PublicId, &user.Login, &user.Password, &user.FirstName, &user.SecondName); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *storage) SelectUserIdByLoginAndPassword(login, password string) (int, error) {
	query := `
		SELECT 
			id
		FROM 
			users 
		WHERE 
			login = $1 and password = $2`

	rows := s.db.QueryRow(s.ctx, query, login, password)

	var id int
	if err := rows.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) SelectUserIdByTocken(tockenText string) (int, error) {
	query := `
		SELECT 
			user_id
		FROM 
			tockens 
		WHERE 
			tocken = $1`

	rows := s.db.QueryRow(s.ctx, query, tockenText)

	var userId int
	if err := rows.Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}

func (s *storage) UpdateUser(u entity.User) error {
	query := `
		UPDATE users SET
			public_id = $2,
			login = $3,
			password = $4,
			first_name = $5,
			second_name = $6
		WHERE
   			id = $1`

	_, err := s.db.Exec(s.ctx, query, u.Id, u.PublicId, u.Login, u.Password, u.FirstName, u.SecondName)
	if err != nil {
		return err
	}

	return err
}

func (s *storage) CreateTocken(userId int, tockenText string) (entity.Tocken, error) {
	query := `
		INSERT INTO
			tockens (user_id, tocken)
		VALUES
			($1, $2)
		RETURNING 
			id`

	tocken := entity.Tocken{
		UserId: userId,
		Text:   tockenText,
	}

	err := s.db.QueryRow(s.ctx, query, tocken.UserId, tocken.Text).Scan(&tocken.Id)
	if err != nil {
		return entity.Tocken{}, err
	}

	return tocken, nil
}

func (s *storage) DeleteTocken(id int) error {
	query := "DELETE FROM tockens WHERE id = $1"

	if _, err := s.db.Exec(s.ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (s *storage) CreateTableUserChats(userId int) error {
	query := `
		CREATE TABLE IF NOT EXISTS chats_user_` + strconv.Itoa(userId) + `(
    		id SERIAL PRIMARY KEY UNIQUE NOT NULL,
			chat_id INTEGER NOT NULL REFERENCES chats (Id) ON DELETE CASCADE
		)`

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) DeleteTableUserChats(userId int) error {
	query := `DROP TABLE IF EXISTS chats_user_` + strconv.Itoa(userId)

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) CreateChat(chat entity.Chat) (int, error) {
	query := `
		INSERT INTO
    		chats (name, description, open)
		VALUES
			($1, $2, $3)
		RETURNING
			id`

	var id int
	err := s.db.QueryRow(s.ctx, query, chat.Description, chat.Name, chat.Open).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) DeleteChat(id int) error {
	query := "DELETE FROM tockens WHERE id = $1"

	if _, err := s.db.Exec(s.ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (s *storage) CreateMessage(chatId int, chat entity.Message) (int, error) {
	query := `
		INSERT INTO
    		messages_chat_` + strconv.Itoa(chatId) + `
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

func (s *storage) DeleteMessage(chatId int, id int) error {
	query := `DELETE FROM messages_chat_` + strconv.Itoa(chatId) + ` WHERE id = $1`

	if _, err := s.db.Exec(s.ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (s *storage) CreateParticipant(chatId int, chat entity.Chat) (int, error) {
	query := `
		INSERT INTO
    		participants_chat_` + strconv.Itoa(chatId) + `
			(user_id, write, post, comment, delete, add_participant)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			id`

	var id int
	err := s.db.QueryRow(s.ctx, query, chat.Description, chat.Name, chat.Open).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) SelectTopParticipants(chatId int, limit int) ([]entity.Participant, error) {
	query := `
		SELECT 
			id, user_id, write, post, comment, delete, add_participant 
		FROM 
			participants_chat_` + strconv.Itoa(chatId) + ` 
		ORDER BY 
			id DESC
		LIMIT $1`

	rows, err := s.db.Query(s.ctx, query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	couriers := make([]entity.Participant, 0, limit)
	for rows.Next() {
		var c entity.Participant

		if err = rows.Scan(&c.Id, &c.UserId, &c.Write, &c.Post, &c.Comment, &c.Delete, &c.AddParticipant); err != nil {
			return nil, err
		}

		couriers = append(couriers, c)
	}

	return couriers, nil
}

func (s *storage) DeleteParticipant(chatId int, id int) error {
	query := `DELETE FROM participants_chat_` + strconv.Itoa(chatId) + ` WHERE id = $1`

	if _, err := s.db.Exec(s.ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (s *storage) CreateTableChatParticipants(chatId int) error {
	query := `
		CREATE TABLE IF NOT EXISTS participants_chat_` + strconv.Itoa(chatId) + `(
    		id SERIAL PRIMARY KEY UNIQUE NOT NULL,
			user_id INTEGER NOT NULL REFERENCES users (Id) ON DELETE CASCADE,
			
			write BOOLEAN NOT NULL,
			post BOOLEAN NOT NULL,
			comment BOOLEAN NOT NULL,
			delete BOOLEAN NOT NULL,
			add_participant BOOLEAN NOT NULL
		)`

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) DeleteTableChatParticipants(chatId int) error {
	query := `DROP TABLE IF EXISTS participants_chat_` + strconv.Itoa(chatId)

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) CreateTableChatMessages(chatId int) error {
	query := `
		CREATE TABLE IF NOT EXISTS messages_chat_` + strconv.Itoa(chatId) + `(
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

func (s *storage) SelectTopMessages(chatId, limit int) ([]entity.Message, error) {
	query := `
		SELECT 
			id, user_id, commented_message_id, message, changed, created 
		FROM 
			messages_chat_` + strconv.Itoa(chatId) + ` 
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
			messages_chat_` + strconv.Itoa(chatId) + `
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

func (s *storage) DeleteTableChatMessages(chatId int) error {
	query := `DROP TABLE IF EXISTS messages_chat_` + strconv.Itoa(chatId)

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}
