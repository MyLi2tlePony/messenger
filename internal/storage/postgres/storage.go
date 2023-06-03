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

func (s *storage) CreateTableParticipantsChat(chatId int) error {
	query := `
		CREATE TABLE IF NOT EXISTS participants_chat_` + strconv.Itoa(chatId) + `(
    		id SERIAL PRIMARY KEY UNIQUE NOT NULL,
			user_id INTEGER NOT NULL REFERENCES users (Id) ON DELETE CASCADE,
			
			write BOOLEAN NOT NULL,
			post BOOLEAN NOT NULL,
			comment BOOLEAN NOT NULL,
			delete BOOLEAN NOT NULL,
			add BOOLEAN NOT NULL
		)`

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) DeleteTableParticipantsChat(chatId int) error {
	query := `DROP TABLE IF EXISTS participants_chat_` + strconv.Itoa(chatId)

	_, err := s.db.Exec(s.ctx, query)
	if err != nil {
		return err
	}

	return nil
}

//// CreateCourier создает запись в базе данных и возвращает id новой записи
//func (s *storage) CreateCourier(ctx context.Context, courier entity.Courier) (id int, err error) {
//	query := `
//		INSERT INTO
//		    Couriers (courier_type, regions, working_hours)
//		VALUES
//		    ($1, $2, $3)
//		RETURNING id`
//
//	err = s.db.QueryRow(ctx, query, courier.CourierType, courier.Regions, courier.WorkingHours).Scan(&id)
//	if err != nil {
//		return 0, err
//	}
//
//	return id, nil
//}
//
//// GetCourierById возвращает данные о курьере
//func (s *storage) GetCourierById(ctx context.Context, id int) (entity.Courier, error) {
//	query := "SELECT id, courier_type, regions, working_hours FROM couriers WHERE id = $1"
//	c := entity.Courier{}
//
//	if err := s.db.QueryRow(ctx, query, id).Scan(&c.Id, &c.CourierType, &c.Regions, &c.WorkingHours); err != nil {
//		return entity.Courier{}, err
//	}
//
//	return c, nil
//}
//
//// DeleteCourierById удаляет курьера с указанным id
//func (s *storage) DeleteCourierById(ctx context.Context, id int) error {
//	query := "DELETE FROM couriers WHERE id = $1"
//
//	if _, err := s.db.Exec(ctx, query, id); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// GetCouriers возвращает срез курьеров с указанными offset и limit
//func (s *storage) GetCouriers(ctx context.Context, offset, limit int) ([]entity.Courier, error) {
//	query := "SELECT id, courier_type, regions, working_hours FROM couriers OFFSET $1 LIMIT $2"
//
//	rows, err := s.db.Query(ctx, query, offset, limit)
//	if err != nil {
//		return nil, err
//	}
//
//	defer rows.Close()
//
//	couriers := make([]entity.Courier, 0, limit)
//	for rows.Next() {
//		var c entity.Courier
//
//		if err = rows.Scan(&c.Id, &c.CourierType, &c.Regions, &c.WorkingHours); err != nil {
//			return nil, err
//		}
//
//		couriers = append(couriers, c)
//	}
//
//	return couriers, nil
//}
//
//// CreateOrder создает запись в базе данных и возвращает id новой записи
//func (s *storage) CreateOrder(ctx context.Context, o entity.Order) (id int, err error) {
//	query := `
//		INSERT INTO
//			orders (weight, regions, delivery_hours, cost, completed_time, courier_id)
//		VALUES
//			($1, $2, $3, $4, $5, $6)
//		RETURNING id`
//
//	err = s.db.QueryRow(ctx, query, o.Weight, o.Regions, o.DeliveryHours, o.Cost, o.CompletedTime, o.CourierId).Scan(&id)
//	if err != nil {
//		return 0, err
//	}
//
//	return id, nil
//}
//
//// GetOrderById возвращает заказ по его id
//func (s *storage) GetOrderById(ctx context.Context, id int) (entity.Order, error) {
//	query := "SELECT id, weight, regions, delivery_hours, cost, completed_time, courier_id FROM orders WHERE id = $1"
//	o := entity.Order{}
//
//	if err := s.db.QueryRow(ctx, query, id).Scan(&o.Id, &o.Weight, &o.Regions, &o.DeliveryHours, &o.Cost, &o.CompletedTime, &o.CourierId); err != nil {
//		return entity.Order{}, err
//	}
//
//	return o, nil
//}
//
//// DeleteOrderById удаляет заказ с указанным id
//func (s *storage) DeleteOrderById(ctx context.Context, id int) error {
//	query := "DELETE FROM orders WHERE id = $1"
//
//	if _, err := s.db.Exec(ctx, query, id); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// GetOrders возвращает срез заказов с указанными offset и limit
//func (s *storage) GetOrders(ctx context.Context, offset, limit int) ([]entity.Order, error) {
//	query := "SELECT id, weight, regions, delivery_hours, cost, completed_time, courier_id FROM orders OFFSET $1 LIMIT $2"
//
//	rows, err := s.db.Query(ctx, query, offset, limit)
//	if err != nil {
//		return nil, err
//	}
//
//	defer rows.Close()
//
//	orders := make([]entity.Order, 0, limit)
//	for rows.Next() {
//		var o entity.Order
//
//		if err = rows.Scan(&o.Id, &o.Weight, &o.Regions, &o.DeliveryHours, &o.Cost, &o.CompletedTime, &o.CourierId); err != nil {
//			return nil, err
//		}
//
//		orders = append(orders, o)
//	}
//
//	return orders, nil
//}
//
//// UpdateOrder изменяет запись
//func (s *storage) UpdateOrder(ctx context.Context, o entity.Order) error {
//	query := `
//		UPDATE orders SET
//		   weight = $2,
//		   regions = $3,
//		   delivery_hours = $4,
//		   cost = $5,
//		   completed_time = $6,
//		   courier_id = $7
//		WHERE id = $1`
//
//	_, err := s.db.Exec(ctx, query, &o.Id, &o.Weight, &o.Regions, &o.DeliveryHours, &o.Cost, &o.CompletedTime, &o.CourierId)
//	if err != nil {
//		return err
//	}
//
//	return err
//}
