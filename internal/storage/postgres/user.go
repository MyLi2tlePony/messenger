package postgres

import (
	"github.com/MyLi2tlePony/messenger/internal/storage/entity"
	"time"
)

// CreateUser создает пользователя с указанным логином и паролем и возвращает присвоенный идентификатор
func (s *storage) CreateUser(login, password string, created time.Time) (int, error) {
	query := `
	INSERT INTO
		users (login, password, first_name, second_name, created)
	VALUES
		($1, $2, '', '', $3)
	RETURNING
		id`

	var id int
	err := s.db.QueryRow(s.ctx, query, login, password, created).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SelectUserByPublicId возвращает данные пользователя с указанным публичным идентификатором
func (s *storage) SelectUserByPublicId(publicId string) (entity.User, error) {
	query := "SELECT id, login, password, first_name, second_name, created FROM users WHERE public_id = $1"

	row := s.db.QueryRow(s.ctx, query, publicId)
	user := entity.User{
		PublicId: publicId,
	}

	if err := row.Scan(&user.Id, &user.Login, &user.Password, &user.FirstName, &user.SecondName, &user.Created); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// SelectUserByToken возвращает данные пользователя с указанным токеном
func (s *storage) SelectUserByToken(token entity.Token) (entity.User, error) {
	query := `
	SELECT 
		id, public_id, login, password, first_name, second_name, created
	FROM 
		users 
	WHERE 
		id = (SELECT user_Id FROM tokens WHERE token = $1)`

	row := s.db.QueryRow(s.ctx, query, token.Text)

	var user entity.User
	if err := row.Scan(&user.Id, &user.PublicId, &user.Login, &user.Password, &user.FirstName, &user.SecondName, &user.Created); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// SelectUserIdByLoginAndPassword возвращает данные пользователя с указанным логином и паролем
func (s *storage) SelectUserIdByLoginAndPassword(login, password string) (int, error) {
	query := `
	SELECT 
		id
	FROM 
		users 
	WHERE 
		login = $1 and password = $2`

	row := s.db.QueryRow(s.ctx, query, login, password)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// SelectUserIdByToken возвращает идентификатор пользователя с указанным токеном
func (s *storage) SelectUserIdByToken(tokenText string) (int, error) {
	query := `
	SELECT 
		user_id
	FROM 
		tokens 
	WHERE 
		token = $1`

	rows := s.db.QueryRow(s.ctx, query, tokenText)

	var userId int
	if err := rows.Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}

// UpdateUser обновляет информацию о пользователе
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

// DeleteUser удаляет пользователя с указанным идентификатором
func (s *storage) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = $1"

	if _, err := s.db.Exec(s.ctx, query, id); err != nil {
		return err
	}

	return nil
}
