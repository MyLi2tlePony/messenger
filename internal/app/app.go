package app

import (
	"errors"
	"github.com/MyLi2tlePony/messenger/internal/server/http/dto"
	"github.com/MyLi2tlePony/messenger/internal/storage/entity"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type Storage interface {
	CreateChat(chat entity.Chat) (int, error)
	SelectChatById(id int) (entity.Chat, error)
	DeleteChat(id int) error

	CreateMessage(chatId int, chat entity.Message) (int, error)
	SelectTopMessages(chatId, limit int) ([]entity.Message, error)
	SelectMessagesById(chatId, minId, maxId int) ([]entity.Message, error)
	DeleteMessage(chatId int, id int) error
	CreateTableChatMessages(chatId int) error
	DeleteTableChatMessages(chatId int) error

	CreateParticipant(chatId int, p entity.Participant) (int, error)
	SelectTopParticipants(chatId int, limit int) ([]entity.Participant, error)
	SelectParticipants(chatId int) ([]entity.Participant, error)
	SelectParticipantIdByUserId(chatId int, userId int) (int, error)
	DeleteParticipant(chatId int, id int) error
	CreateTableChatParticipants(chatId int) error
	DeleteTableChatParticipants(chatId int) error

	CreateToken(userId int, tokenText string) (int, error)
	DeleteToken(id int) error

	CreateUser(login, password string, created time.Time) (int, error)
	SelectUserByPublicId(publicId string) (entity.User, error)
	SelectUserByToken(token entity.Token) (entity.User, error)
	SelectUserIdByLoginAndPassword(login, password string) (int, error)
	SelectUserIdByToken(tokenText string) (int, error)
	UpdateUser(u entity.User) error
	DeleteUser(id int) error

	CreateUserChat(userId, chatId int) (int, error)
	SelectUserChats(userId int) ([]entity.UserChats, error)
	DeleteUserChat(userId, userChatId int) error
	CreateTableUserChats(userId int) error
	DeleteTableUserChats(userId int) error
}

type app struct {
	storage Storage
}

var (
	ErrAccessDenied = errors.New("err access denied")
	ErrNotFound     = errors.New("err not found")
)

func New(storage Storage) *app {
	return &app{
		storage: storage,
	}
}

func (a *app) CreateUser(request dto.CreateUserRequest) error {
	id, err := a.storage.CreateUser(request.Login, request.Password, time.Now().Round(time.Second).UTC())
	if err != nil {
		return err
	}

	err = a.storage.CreateTableUserChats(id)
	if err != nil {
		return err
	}

	err = a.storage.UpdateUser(entity.User{
		Id:         id,
		Login:      request.Login,
		Password:   request.Password,
		PublicId:   strconv.Itoa(id),
		FirstName:  strconv.Itoa(id),
		SecondName: strconv.Itoa(id),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *app) CreateToken(request dto.CreateTokenRequest) (dto.Token, error) {
	userId, err := a.storage.SelectUserIdByLoginAndPassword(request.Login, request.Password)
	if err != nil {
		return dto.Token{}, err
	}

	dtoToken := dto.Token{
		Text: uuid.New().String(),
	}

	_, err = a.storage.CreateToken(userId, dtoToken.Text)
	if err != nil {
		return dto.Token{}, err
	}

	return dtoToken, nil
}

func (a *app) SelectUserByPublicId(publicId string) (dto.User, error) {
	user, err := a.storage.SelectUserByPublicId(publicId)
	if err != nil {
		return dto.User{}, err
	}

	dtoUser := dto.User{
		PublicId:   user.PublicId,
		Login:      user.Login,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Created:    user.Created,
	}

	return dtoUser, nil
}

func (a *app) SelectUserByToken(dtoToken dto.Token) (dto.User, error) {
	token := entity.Token{
		Text: dtoToken.Text,
	}

	user, err := a.storage.SelectUserByToken(token)
	if err != nil {
		return dto.User{}, err
	}

	dtoUser := dto.User{
		PublicId:   user.PublicId,
		Login:      user.Login,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Created:    user.Created,
	}

	return dtoUser, nil
}

func (a *app) UpdateUser(dtoToken dto.Token, dtoUser dto.User) error {
	token := entity.Token{
		Text: dtoToken.Text,
	}

	userId, err := a.storage.SelectUserIdByToken(token.Text)
	if err != nil {
		return err
	}

	user := entity.User{
		Id:         userId,
		PublicId:   dtoUser.PublicId,
		Login:      dtoUser.Login,
		FirstName:  dtoUser.FirstName,
		SecondName: dtoUser.SecondName,
		Created:    dtoUser.Created,
	}

	err = a.storage.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

//
////////////////////////////////////////////////////////////////////////////////////////////////////////////
//

func (a *app) CreateChat(createChatRequest dto.CreateChatRequest) error {
	chat := entity.Chat{
		Name:        createChatRequest.Chat.Name,
		Description: createChatRequest.Chat.Description,
		Open:        createChatRequest.Chat.Open,
		Created:     createChatRequest.Chat.Created,
	}
	chatId, err := a.storage.CreateChat(chat)
	if err != nil {
		return err
	}

	err = a.storage.CreateTableChatMessages(chatId)
	if err != nil {
		return err
	}

	err = a.storage.CreateTableChatParticipants(chatId)
	if err != nil {
		return err
	}

	for _, p := range createChatRequest.Participants {
		participant := entity.Participant{
			UserId: p.UserId,

			Write:   p.Write,
			Post:    p.Post,
			Comment: p.Comment,
			Delete:  p.Delete,

			AddParticipant: p.AddParticipant,
		}

		if _, err = a.storage.CreateParticipant(chatId, participant); err != nil {
			return err
		}

		if _, err = a.storage.CreateUserChat(participant.UserId, chatId); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) SelectChatById(id int) (dto.Chat, error) {
	chat, err := a.storage.SelectChatById(id)
	if err != nil {
		return dto.Chat{}, err
	}

	dtoChat := dto.Chat{
		Id: chat.Id,

		Description: chat.Description,
		Name:        chat.Name,

		Open:    chat.Open,
		Created: chat.Created,
	}

	return dtoChat, nil
}

func (a *app) CreateMessage(chatId int, dtoToken dto.Token, dtoMessage dto.Message) error {
	token := entity.Token{
		Text: dtoToken.Text,
	}
	userId, err := a.storage.SelectUserIdByToken(token.Text)
	if err != nil {
		return err
	}

	userChats, err := a.storage.SelectUserChats(userId)
	if err != nil {
		return err
	}

	message := entity.Message{
		Id:     dtoMessage.Id,
		UserId: dtoMessage.Id,

		Changed: dtoMessage.Changed,
		Text:    dtoMessage.Text,

		CommentedMessageId: dtoMessage.CommentedMessageId,

		Created: dtoMessage.Created,
	}

	for i := range userChats {
		if userChats[i].ChatId == chatId {
			if _, err = a.storage.CreateMessage(chatId, message); err != nil {
				return err
			}

			return nil
		}
	}

	return ErrAccessDenied
}

func (a *app) SelectTopMessages(dtoToken dto.Token, chatId, limit int) ([]dto.Message, error) {
	chat, err := a.storage.SelectChatById(chatId)
	if err != nil {
		return nil, err
	}

	if chat.Open {
		messages, err := a.storage.SelectTopMessages(chatId, limit)
		if err != nil {
			return nil, err
		}

		dtoMessages := make([]dto.Message, len(messages))
		for j, m := range messages {
			dtoMessages[j] = dto.Message{
				Id:     m.Id,
				UserId: m.Id,

				Changed: m.Changed,
				Text:    m.Text,

				CommentedMessageId: m.CommentedMessageId,

				Created: m.Created,
			}

		}

		return dtoMessages, nil
	}

	token := entity.Token{
		Text: dtoToken.Text,
	}
	userId, err := a.storage.SelectUserIdByToken(token.Text)
	if err != nil {
		return nil, err
	}

	userChats, err := a.storage.SelectUserChats(userId)
	if err != nil {
		return nil, err
	}

	for i := range userChats {
		if userChats[i].ChatId == chatId {
			messages, err := a.storage.SelectTopMessages(chatId, limit)
			if err != nil {
				return nil, err
			}

			dtoMessages := make([]dto.Message, len(messages))
			for j, m := range messages {
				dtoMessages[j] = dto.Message{
					Id:     m.Id,
					UserId: m.Id,

					Changed: m.Changed,
					Text:    m.Text,

					CommentedMessageId: m.CommentedMessageId,

					Created: m.Created,
				}

			}

			return dtoMessages, nil
		}
	}

	return nil, ErrAccessDenied
}

func (a *app) SelectMessagesById(chatId, minId, maxId int, dtoToken dto.Token) ([]dto.Message, error) {
	chat, err := a.storage.SelectChatById(chatId)
	if err != nil {
		return nil, err
	}

	if chat.Open {
		messages, err := a.storage.SelectMessagesById(chatId, minId, maxId)
		if err != nil {
			return nil, err
		}

		dtoMessages := make([]dto.Message, len(messages))
		for j, m := range messages {
			dtoMessages[j] = dto.Message{
				Id:     m.Id,
				UserId: m.Id,

				Changed: m.Changed,
				Text:    m.Text,

				CommentedMessageId: m.CommentedMessageId,

				Created: m.Created,
			}

		}

		return dtoMessages, nil
	}

	token := entity.Token{
		Text: dtoToken.Text,
	}
	userId, err := a.storage.SelectUserIdByToken(token.Text)
	if err != nil {
		return nil, err
	}

	userChats, err := a.storage.SelectUserChats(userId)
	if err != nil {
		return nil, err
	}

	for i := range userChats {
		if userChats[i].ChatId == chatId {
			messages, err := a.storage.SelectMessagesById(chatId, minId, maxId)
			if err != nil {
				return nil, err
			}

			dtoMessages := make([]dto.Message, len(messages))
			for j, m := range messages {
				dtoMessages[j] = dto.Message{
					Id:     m.Id,
					UserId: m.Id,

					Changed: m.Changed,
					Text:    m.Text,

					CommentedMessageId: m.CommentedMessageId,

					Created: m.Created,
				}

			}

			return dtoMessages, nil
		}
	}

	return nil, ErrAccessDenied
}

func (a *app) DeleteMessage(chatId, messageId int, dtoToken dto.Token) error {
	token := entity.Token{
		Text: dtoToken.Text,
	}
	userId, err := a.storage.SelectUserIdByToken(token.Text)
	if err != nil {
		return err
	}

	userChats, err := a.storage.SelectUserChats(userId)
	if err != nil {
		return err
	}

	for i := range userChats {
		if userChats[i].ChatId == chatId {
			if err = a.storage.DeleteMessage(chatId, messageId); err != nil {
				return err
			}

			return nil
		}
	}

	return ErrAccessDenied
}

func (a *app) GetUserChats(dtoToken dto.Token) ([]dto.Chat, error) {
	token := entity.Token{
		Text: dtoToken.Text,
	}
	userId, err := a.storage.SelectUserIdByToken(token.Text)
	if err != nil {
		return nil, err
	}

	userChats, err := a.storage.SelectUserChats(userId)
	if err != nil {
		return nil, err
	}

	dtoChats := make([]dto.Chat, len(userChats))
	for i, c := range userChats {
		chat, err := a.storage.SelectChatById(c.ChatId)
		if err != nil {
			return nil, err
		}

		dtoChats[i] = dto.Chat{
			Id: chat.Id,

			Description: chat.Description,
			Name:        chat.Name,

			Open:    chat.Open,
			Created: chat.Created,
		}
	}

	return dtoChats, nil
}

func (a *app) DeleteUserChat(dtoToken dto.Token, chatId int) error {
	token := entity.Token{
		Text: dtoToken.Text,
	}
	userId, err := a.storage.SelectUserIdByToken(token.Text)
	if err != nil {
		return err
	}

	userChats, err := a.storage.SelectUserChats(userId)
	if err != nil {
		return err
	}

	for i := range userChats {
		if userChats[i].ChatId == chatId {
			if err = a.storage.DeleteUser(chatId); err != nil {
				return err
			}

			id, err := a.storage.SelectParticipantIdByUserId(chatId, userId)
			if err != nil {
				return err
			}

			if err = a.storage.DeleteParticipant(chatId, id); err != nil {
				return err
			}

			return nil
		}
	}

	return ErrAccessDenied
}

func (a *app) CreateParticipant(dtoToken dto.Token, chatId int, dtoParticipant dto.Participant) error {
	token := entity.Token{
		Text: dtoToken.Text,
	}
	userId, err := a.storage.SelectUserIdByToken(token.Text)
	if err != nil {
		return err
	}

	participants, err := a.storage.SelectParticipants(chatId)
	if err != nil {
		return err
	}

	for i := range participants {
		if participants[i].UserId == userId {
			if !participants[i].AddParticipant {
				break
			}

			participant := entity.Participant{
				UserId:            dtoParticipant.UserId,
				Write:             dtoParticipant.Write,
				Post:              dtoParticipant.Post,
				Comment:           dtoParticipant.Comment,
				Delete:            dtoParticipant.Delete,
				AddParticipant:    dtoParticipant.AddParticipant,
				DeleteParticipant: dtoParticipant.DeleteParticipant,
			}

			if _, err = a.storage.CreateParticipant(chatId, participant); err != nil {
				return err
			}

			return nil
		}
	}

	return ErrAccessDenied
}

func (a *app) SelectParticipant(chatId int) ([]dto.Participant, error) {
	participants, err := a.storage.SelectParticipants(chatId)
	if err != nil {
		return nil, err
	}

	dtoParticipants := make([]dto.Participant, len(participants))
	for i, p := range participants {
		dtoParticipants[i] = dto.Participant{
			UserId:            p.UserId,
			Write:             p.Write,
			Post:              p.Post,
			Comment:           p.Comment,
			Delete:            p.Delete,
			AddParticipant:    p.AddParticipant,
			DeleteParticipant: p.DeleteParticipant,
		}
	}

	return dtoParticipants, nil
}

func (a *app) DeleteParticipant(dtoToken dto.Token, chatId, participantId int) error {
	token := entity.Token{
		Text: dtoToken.Text,
	}
	userId, err := a.storage.SelectUserIdByToken(token.Text)
	if err != nil {
		return err
	}

	participants, err := a.storage.SelectParticipants(chatId)
	if err != nil {
		return err
	}

	for i := range participants {
		if participants[i].UserId == userId {
			if !participants[i].DeleteParticipant {
				return ErrAccessDenied
			}

			err = a.storage.DeleteParticipant(chatId, participantId)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return ErrAccessDenied
}
