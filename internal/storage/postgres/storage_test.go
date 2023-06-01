///go:build integration

package postgres

import (
	"context"
	"github.com/google/uuid"
	"path"
	"testing"

	databaseConfig "github.com/MyLi2tlePony/messenger/internal/config/database"
	"github.com/MyLi2tlePony/messenger/internal/storage/entity"

	"github.com/stretchr/testify/require"
)

var configPath = path.Join("..", "..", "..", "configs", "test", "config.toml")

func TestUser(t *testing.T) {
	dbConfig, err := databaseConfig.New(configPath)
	require.Nil(t, err)

	ctx := context.Background()
	db, err := Connect(ctx, dbConfig.GetConnectionString())
	require.Nil(t, err)

	users := [3]entity.User{
		{PublicId: "AndreyId", Login: "Andrey", Password: "1234", FirstName: "Andreyka", SecondName: "Anasd"},
		{PublicId: "DimaId", Login: "Dima", Password: "12312344", FirstName: "Dimaka", SecondName: "Dimaasdf"},
		{PublicId: "AsadfId", Login: "asdf", Password: "1234", FirstName: "ASDfka", SecondName: "Aasdf"},
	}
	tockens := [3]entity.Tocken{}

	for i := range users {
		users[i].Id, err = db.CreateUser(users[i].Login, users[i].Password)
		require.Nil(t, err)

		require.Nil(t, db.UpdateUser(users[i]))

		tockens[i], err = db.CreateTocken(users[i].Id, uuid.New().String())
		require.Nil(t, err)
	}

	for i, u := range users {
		var user entity.User

		user, err = db.SelectUserByPublicId(u.PublicId)
		require.Nil(t, err)
		require.True(t, user.Equals(users[i]))

		userId, err := db.SelectUserIdByLoginAndPassword(u.Login, u.Password)
		require.Nil(t, err)
		require.True(t, users[i].Id == userId)

		user, err = db.SelectUserByTocken(tockens[i])
		require.Nil(t, err)
		require.True(t, user.Equals(users[i]))
	}

	for _, tocken := range tockens {
		require.Nil(t, db.DeleteTocken(tocken.Id))
	}

	for _, user := range users {
		require.Nil(t, db.DeleteUser(user.Id))
	}
}

//func TestCouriers(t *testing.T) {
//	dbConfig, err := databaseConfig.New(configPath)
//	require.Nil(t, err)
//
//	ctx := context.Background()
//	db, err := Connect(ctx, dbConfig.GetConnectionString())
//	require.Nil(t, err)
//
//	couriers := []*entity.Courier{
//		entity.NewCourier("FOOT", []int{1, 2}, []string{"12:15-17:15"}),
//		entity.NewCourier("FOOT", []int{3, 4}, []string{"13:15-18:15"}),
//		entity.NewCourier("BIKE", []int{5, 6}, []string{"14:15-19:15"}),
//		entity.NewCourier("BIKE", []int{7, 8}, []string{"15:15-20:15"}),
//		entity.NewCourier("AUTO", []int{9, 10}, []string{"16:15-21:15"}),
//		entity.NewCourier("AUTO", []int{11, 12}, []string{"17:15-22:15"}),
//	}
//
//	for i := range couriers {
//		id, err := db.CreateCourier(ctx, *couriers[i])
//		require.Nil(t, err)
//
//		couriers[i].Id = id
//	}
//
//	for i := range couriers {
//		courier, err := db.GetCourierById(ctx, couriers[i].Id)
//		require.Nil(t, err)
//
//		require.True(t, entity.CouriersEquals(courier, *couriers[i]))
//	}
//
//	offset, limit := 0, 2
//	couriersPart, err := db.GetCouriers(ctx, offset, limit)
//	require.Nil(t, err)
//
//	for i := 0; i+offset < offset+limit; i++ {
//		require.True(t, entity.CouriersEquals(*couriers[i+offset], couriersPart[i]))
//	}
//
//	offset, limit = 2, 4
//	couriersPart, err = db.GetCouriers(ctx, offset, limit)
//	require.Nil(t, err)
//
//	for i := 0; i+offset < offset+limit; i++ {
//		require.True(t, entity.CouriersEquals(*couriers[i+offset], couriersPart[i]))
//	}
//
//	for i := range couriers {
//		require.Nil(t, db.DeleteCourierById(ctx, couriers[i].Id))
//	}
//}
//
//func TestOrders(t *testing.T) {
//	dbConfig, err := databaseConfig.New(configPath)
//	require.Nil(t, err)
//
//	ctx := context.Background()
//	db, err := Connect(ctx, dbConfig.GetConnectionString())
//	require.Nil(t, err)
//
//	orders := []*entity.Order{
//		entity.NewOrder(10, 1, []string{"12:15-17:15"}, 100),
//		entity.NewOrder(20, 2, []string{"12:15-17:15"}, 111),
//		entity.NewOrder(30, 2, []string{"12:15-17:15"}, 222),
//		entity.NewOrder(40, 3, []string{"12:15-17:15"}, 333),
//		entity.NewOrder(11, 4, []string{"12:15-17:15"}, 444),
//		entity.NewOrder(15, 5, []string{"12:15-17:15"}, 555),
//	}
//
//	for i := range orders {
//		id, err := db.CreateOrder(ctx, *orders[i])
//		require.Nil(t, err)
//
//		orders[i].Id = id
//	}
//
//	for i := range orders {
//		courier, err := db.GetOrderById(ctx, orders[i].Id)
//		require.Nil(t, err)
//
//		require.True(t, entity.OrderEquals(courier, *orders[i]))
//	}
//
//	offset, limit := 0, 2
//	ordersPart, err := db.GetOrders(ctx, offset, limit)
//	require.Nil(t, err)
//
//	for i := 0; i+offset < offset+limit; i++ {
//		require.True(t, entity.OrderEquals(*orders[i+offset], ordersPart[i]))
//	}
//
//	offset, limit = 2, 4
//	ordersPart, err = db.GetOrders(ctx, offset, limit)
//	require.Nil(t, err)
//
//	for i := 0; i+offset < offset+limit; i++ {
//		require.True(t, entity.OrderEquals(*orders[i+offset], ordersPart[i]))
//	}
//
//	for i := range orders {
//		orders[i].CompletedTime = strconv.Itoa(i + 1)
//		orders[i].CourierId = i + 1
//
//		err = db.UpdateOrder(ctx, *orders[i])
//		require.Nil(t, err)
//	}
//
//	for i := range orders {
//		courier, err := db.GetOrderById(ctx, orders[i].Id)
//		require.Nil(t, err)
//
//		require.True(t, entity.OrderEquals(courier, *orders[i]))
//	}
//
//	for i := range orders {
//		require.Nil(t, db.DeleteOrderById(ctx, orders[i].Id))
//	}
//}
