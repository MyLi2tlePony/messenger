///go:build integration

package postgres

import (
	"context"
	"github.com/google/uuid"
	"path"
	"strconv"
	"testing"
	"time"

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

		err = db.CreateTableUserChats(users[i].Id)
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

	chats := []entity.Chat{
		{Name: "Новый чат 1", Description: "Описание чата 1", Open: true, Created: time.Now()},
		{Name: "Новый чат 2", Description: "Описание чата 2", Open: false, Created: time.Now()},
		{Name: "Новый чат 3", Description: "Описание чата 3", Open: true, Created: time.Now()},
	}

	for i := range chats {
		chats[i].Id, err = db.CreateChat(chats[i])
		require.Nil(t, err)

		require.Nil(t, db.CreateTableChatParticipants(chats[i].Id))
		require.Nil(t, db.CreateTableChatMessages(chats[i].Id))

		messages := make([]entity.Message, 0)

		for j, user := range users {
			messages = append(messages, entity.Message{
				UserId:  user.Id,
				Changed: false,
				Text:    strconv.Itoa(j),
				Created: time.Now()})

			messages[i].Id, err = db.CreateMessage(chats[i].Id, messages[j])

		}

	}

	for _, chat := range chats {
		require.Nil(t, db.DeleteChat(chat.Id))
		require.Nil(t, db.DeleteTableChatParticipants(chat.Id))
		require.Nil(t, db.DeleteTableChatMessages(chat.Id))
	}

	for _, tocken := range tockens {
		require.Nil(t, db.DeleteTocken(tocken.Id))
	}

	for _, user := range users {
		require.Nil(t, db.DeleteUser(user.Id))
		require.Nil(t, db.DeleteTableUserChats(user.Id))
	}
}
