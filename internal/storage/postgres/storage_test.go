///go:build integration

package postgres

import (
	"context"
	databaseConfig "github.com/MyLi2tlePony/messenger/internal/config/database"
	"github.com/MyLi2tlePony/messenger/internal/storage/entity"
	"github.com/google/uuid"
	"path"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var configPath = path.Join("..", "..", "..", "configs", "test", "config.toml")

// TestUser тестирует все функции, связанные со структурой пользователя
func TestUser(t *testing.T) {
	dbConfig, err := databaseConfig.New(configPath)
	require.Nil(t, err)

	ctx := context.Background()
	db, err := Connect(ctx, dbConfig.GetConnectionString())
	require.Nil(t, err)

	users := []entity.User{
		entity.NewUser("AndreyId", "Andrey", "1234", "Andreyka", "Anasd", time.Now().Round(time.Second).UTC()),
		entity.NewUser("DimaId", "Dima", "12312344", "Dimaka", "Dimaasdf", time.Now().Round(time.Second).UTC()),
		entity.NewUser("AsadfId", "asdf", "1234", "ASDfka", "Aasdf", time.Now().Round(time.Second).UTC()),
	}

	for i := range users {
		users[i].Id, err = db.CreateUser(users[i].Login, users[i].Password, users[i].Created)
		require.Nil(t, err)

		require.Nil(t, db.UpdateUser(users[i]))
	}

	for i, u := range users {
		var user entity.User

		user, err = db.SelectUserByPublicId(u.PublicId)
		require.Nil(t, err)
		require.Equal(t, users[i], user)

		userId, err := db.SelectUserIdByLoginAndPassword(u.Login, u.Password)
		require.Nil(t, err)
		require.True(t, users[i].Id == userId)
	}

	for _, user := range users {
		require.Nil(t, db.DeleteUser(user.Id))
	}
}

// TestToken тестирует все функции, связанные со структурой токена
func TestToken(t *testing.T) {
	dbConfig, err := databaseConfig.New(configPath)
	require.Nil(t, err)

	ctx := context.Background()
	db, err := Connect(ctx, dbConfig.GetConnectionString())
	require.Nil(t, err)

	tokens := []entity.Token{
		{Text: uuid.New().String()},
		{Text: uuid.New().String()},
		{Text: uuid.New().String()},
	}

	users := []entity.User{
		entity.NewUser("AndreyId", "Andrey", "1234", "Andreyka", "Anasd", time.Now().Round(time.Second).UTC()),
		entity.NewUser("DimaId", "Dima", "12312344", "Dimaka", "Dimaasdf", time.Now().Round(time.Second).UTC()),
		entity.NewUser("AsadfId", "asdf", "1234", "ASDfka", "Aasdf", time.Now().Round(time.Second).UTC()),
	}

	for i := range users {
		users[i].Id, err = db.CreateUser(users[i].Login, users[i].Password, users[i].Created)
		require.Nil(t, err)

		tokens[i].Id, err = db.CreateToken(users[i].Id, tokens[i].Text)
		tokens[i].UserId = users[i].Id
		require.Nil(t, err)

		require.Nil(t, db.UpdateUser(users[i]))
	}

	for i := range users {
		var user entity.User

		user, err = db.SelectUserByToken(tokens[i])
		require.Nil(t, err)
		require.Equal(t, users[i], user)

		userId, err := db.SelectUserIdByToken(tokens[i].Text)
		require.Nil(t, err)
		require.Equal(t, users[i].Id, userId)
	}

	for _, token := range tokens {
		require.Nil(t, db.DeleteToken(token.Id))
	}

	for _, user := range users {
		require.Nil(t, db.DeleteUser(user.Id))
	}
}

// TestChats тестирует все функции, связанные со структурой чатов
func TestChats(t *testing.T) {
	dbConfig, err := databaseConfig.New(configPath)
	require.Nil(t, err)

	ctx := context.Background()
	db, err := Connect(ctx, dbConfig.GetConnectionString())
	require.Nil(t, err)

	chats := []entity.Chat{
		{Name: "Новый чат 1", Description: "Описание чата 1", Open: true, Created: time.Now().Round(time.Second).UTC()},
		{Name: "Новый чат 2", Description: "Описание чата 2", Open: false, Created: time.Now().Round(time.Second).UTC()},
		{Name: "Новый чат 3", Description: "Описание чата 3", Open: true, Created: time.Now().Round(time.Second).UTC()},
	}

	for i := range chats {
		chats[i].Id, err = db.CreateChat(chats[i])
		require.Nil(t, err)

		chat, err := db.SelectChatById(chats[i].Id)
		require.Nil(t, err)
		require.Equal(t, chats[i], chat)
	}

	for i := range chats {
		require.Nil(t, db.DeleteChat(chats[i].Id))
	}
}

func TestUserChats(t *testing.T) {
	dbConfig, err := databaseConfig.New(configPath)
	require.Nil(t, err)

	ctx := context.Background()
	db, err := Connect(ctx, dbConfig.GetConnectionString())
	require.Nil(t, err)

	chats := []entity.Chat{
		{Name: "Новый чат 1", Description: "Описание чата 1", Open: true, Created: time.Now().Round(time.Second).UTC()},
		{Name: "Новый чат 2", Description: "Описание чата 2", Open: false, Created: time.Now().Round(time.Second).UTC()},
		{Name: "Новый чат 3", Description: "Описание чата 3", Open: true, Created: time.Now().Round(time.Second).UTC()},
	}

	for i := range chats {
		chats[i].Id, err = db.CreateChat(chats[i])
		require.Nil(t, err)
	}

	users := []entity.User{
		entity.NewUser("AndreyId", "Andrey", "1234", "Andreyka", "Anasd", time.Now().Round(time.Second).UTC()),
		entity.NewUser("DimaId", "Dima", "12312344", "Dimaka", "Dimaasdf", time.Now().Round(time.Second).UTC()),
		entity.NewUser("AsadfId", "asdf", "1234", "ASDfka", "Aasdf", time.Now().Round(time.Second).UTC()),
	}

	for i := range users {
		users[i].Id, err = db.CreateUser(users[i].Login, users[i].Password, users[i].Created)
		require.Nil(t, err)

		require.Nil(t, db.CreateTableUserChats(users[i].Id))
	}

	for _, u := range users {
		ids := make([]int, len(chats))

		for i, c := range chats {
			ids[i], err = db.CreateUserChat(u.Id, c.Id)
			require.Nil(t, err)
		}

		selectedUserChats, err := db.SelectUserChats(u.Id)
		require.Nil(t, err)
		require.Equal(t, len(selectedUserChats), len(chats))

		for i := range chats {
			require.Equal(t, chats[i].Id, selectedUserChats[i].ChatId)
		}

		for _, id := range ids {
			require.Nil(t, db.DeleteUserChat(u.Id, id))
		}

		require.Nil(t, db.DeleteTableUserChats(u.Id))
	}

	for _, user := range users {
		require.Nil(t, db.DeleteUser(user.Id))
	}

	for _, chat := range chats {
		require.Nil(t, db.DeleteChat(chat.Id))
	}
}

func TestMessage(t *testing.T) {
	dbConfig, err := databaseConfig.New(configPath)
	require.Nil(t, err)

	ctx := context.Background()
	db, err := Connect(ctx, dbConfig.GetConnectionString())
	require.Nil(t, err)

	users := []entity.User{
		entity.NewUser("AndreyId", "Andrey", "1234", "Andreyka", "Anasd", time.Now().Round(time.Second).UTC()),
		entity.NewUser("DimaId", "Dima", "12312344", "Dimaka", "Dimaasdf", time.Now().Round(time.Second).UTC()),
		entity.NewUser("AsadfId", "asdf", "1234", "ASDfka", "Aasdf", time.Now().Round(time.Second).UTC()),
	}

	for i := range users {
		users[i].Id, err = db.CreateUser(users[i].Login, users[i].Password, users[i].Created)
		require.Nil(t, err)
	}

	chats := []entity.Chat{
		{Name: "Новый чат 1", Description: "Описание чата 1", Open: true, Created: time.Now().Round(time.Second).UTC()},
		{Name: "Новый чат 2", Description: "Описание чата 2", Open: false, Created: time.Now().Round(time.Second).UTC()},
		{Name: "Новый чат 3", Description: "Описание чата 3", Open: true, Created: time.Now().Round(time.Second).UTC()},
	}

	for i := range chats {
		chats[i].Id, err = db.CreateChat(chats[i])
		require.Nil(t, err)

		require.Nil(t, db.CreateTableChatMessages(chats[i].Id))
	}

	for i := range chats {
		messages := make([]entity.Message, 0)

		for j, user := range users {
			messages = append(messages, entity.Message{
				UserId:  user.Id,
				Changed: false,
				Text:    strconv.Itoa(j),
				Created: time.Now().Round(time.Second).UTC(),
			})

			messages[j].Id, err = db.CreateMessage(chats[i].Id, messages[j])
			require.Nil(t, err)
		}

		selectedMessages, err := db.SelectTopMessages(chats[i].Id, len(users))
		require.Nil(t, err)

		sort.Slice(selectedMessages, func(i, j int) bool {
			return selectedMessages[i].Id < selectedMessages[j].Id
		})

		require.True(t, len(selectedMessages) == len(messages))

		for j := range messages {
			require.Equal(t, selectedMessages[j], messages[j])
		}

		selectedMessages, err = db.SelectMessagesById(chats[i].Id, messages[0].Id, messages[len(messages)-1].Id)
		require.Nil(t, err)

		sort.Slice(selectedMessages, func(i, j int) bool {
			return selectedMessages[i].Id < selectedMessages[j].Id
		})

		require.True(t, len(selectedMessages) == len(messages))

		for j := range messages {
			require.Equal(t, selectedMessages[j], messages[j])
		}

		for j := range messages {
			require.Nil(t, db.DeleteMessage(chats[i].Id, messages[j].Id))
		}
	}

	for _, chat := range chats {
		require.Nil(t, db.DeleteChat(chat.Id))
		require.Nil(t, db.DeleteTableChatMessages(chat.Id))
	}

	for _, user := range users {
		require.Nil(t, db.DeleteUser(user.Id))
	}
}

func TestParticipant(t *testing.T) {
	dbConfig, err := databaseConfig.New(configPath)
	require.Nil(t, err)

	ctx := context.Background()
	db, err := Connect(ctx, dbConfig.GetConnectionString())
	require.Nil(t, err)

	users := []entity.User{
		entity.NewUser("AndreyId", "Andrey", "1234", "Andreyka", "Anasd", time.Now().Round(time.Second).UTC()),
		entity.NewUser("DimaId", "Dima", "12312344", "Dimaka", "Dimaasdf", time.Now().Round(time.Second).UTC()),
		entity.NewUser("AsadfId", "asdf", "1234", "ASDfka", "Aasdf", time.Now().Round(time.Second).UTC()),
	}

	for i := range users {
		users[i].Id, err = db.CreateUser(users[i].Login, users[i].Password, users[i].Created)
		require.Nil(t, err)
	}

	chats := []entity.Chat{
		{Name: "Новый чат 1", Description: "Описание чата 1", Open: true, Created: time.Now().Round(time.Second).UTC()},
		{Name: "Новый чат 2", Description: "Описание чата 2", Open: false, Created: time.Now().Round(time.Second).UTC()},
		{Name: "Новый чат 3", Description: "Описание чата 3", Open: true, Created: time.Now().Round(time.Second).UTC()},
	}

	for i := range chats {
		chats[i].Id, err = db.CreateChat(chats[i])
		require.Nil(t, err)

		require.Nil(t, db.CreateTableChatParticipants(chats[i].Id))
	}

	for i := range chats {
		participants := make([]entity.Participant, 0)
		for j, user := range users {
			participants = append(participants, entity.Participant{
				UserId: user.Id,
			})

			participants[j].Id, err = db.CreateParticipant(chats[i].Id, participants[j])
			require.Nil(t, err)

			participantId, err := db.SelectParticipantIdByUserId(chats[i].Id, participants[j].UserId)
			require.Nil(t, err)
			require.Equal(t, participants[j].Id, participantId)
		}

		selectedParticipants, err := db.SelectTopParticipants(chats[i].Id, len(users))
		require.Nil(t, err)

		sort.Slice(selectedParticipants, func(i, j int) bool {
			return selectedParticipants[i].Id < selectedParticipants[j].Id
		})

		require.True(t, len(selectedParticipants) == len(participants))

		for j := range participants {
			require.Equal(t, selectedParticipants[j], participants[j])
		}

		selectedParticipants, err = db.SelectParticipants(chats[i].Id)
		require.Nil(t, err)

		sort.Slice(selectedParticipants, func(i, j int) bool {
			return selectedParticipants[i].Id < selectedParticipants[j].Id
		})

		require.True(t, len(selectedParticipants) == len(participants))

		for j := range participants {
			require.Equal(t, selectedParticipants[j], participants[j])
		}

		for j := range participants {
			require.Nil(t, db.DeleteParticipant(chats[i].Id, participants[j].Id))
		}
	}

	for _, chat := range chats {
		require.Nil(t, db.DeleteChat(chat.Id))
		require.Nil(t, db.DeleteTableChatParticipants(chat.Id))
	}

	for _, user := range users {
		require.Nil(t, db.DeleteUser(user.Id))
	}
}
