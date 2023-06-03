package app

import (
	"github.com/MyLi2tlePony/messenger/internal/server/http/dto"
	"github.com/MyLi2tlePony/messenger/internal/storage/entity"
	"github.com/google/uuid"
)

type Storage interface {
	CreateUser(login, password string) (int, error)
	SelectUserByPublicId(publicId string) (entity.User, error)
	SelectUserByTocken(tocken entity.Tocken) (entity.User, error)
	SelectUserIdByLoginAndPassword(login, password string) (int, error)
	UpdateUser(user entity.User) error
	SelectUserIdByTocken(tockenText string) (int, error)
	CreateTocken(userId int, tockenText string) (entity.Tocken, error)

	CreateTableUserChats(userId int) error

	//CreateChat(tocken dto.Tocken, chatName string, chatType int, adminsIds, userIds []int) error

	//CreateMessage(chatId int, message dto.Message) error
	//GetMessage(tocken dto.Tocken, messageId, chatId int) (dto.Message, error)
	//GetPrevMessages(tocken dto.Tocken, chatId, offset, count int) ([]dto.Message, error)
	//UpdateMessage(tocken dto.Tocken, chatId int, message dto.Message) error
	//SharedMessage(tocken dto.Tocken, srcChatId, dstChatId int, messages []dto.Message) error
	//CreateComment(tocken dto.Tocken, messageId int, message dto.Message) error

	//SelectChatById(id int) (dto.Chat, error)
	//GetPrevChats(tocken dto.Tocken, offset, count int) ([]dto.Chat, error)
	//UpdateChat(tocken dto.Tocken, chat dto.Chat) error
	//
	//CheckUserAcess(chatId, userId int) (bool, error)
}

type app struct {
	storage Storage
}

func New(storage Storage) *app {
	return &app{
		storage: storage,
	}
}

func (a *app) CreateUser(request dto.CreateUserRequest) error {
	id, err := a.storage.CreateUser(request.Login, request.Password)
	if err != nil {
		return err
	}

	err = a.storage.CreateTableUserChats(id)
	if err != nil {
		return err
	}

	return nil
}

func (a *app) CreateTocken(request dto.CreateTockenRequest) (dto.Tocken, error) {
	userId, err := a.storage.SelectUserIdByLoginAndPassword(request.Login, request.Password)
	if err != nil {
		return dto.Tocken{}, err
	}

	tocken, err := a.storage.CreateTocken(userId, uuid.New().String())
	if err != nil {
		return dto.Tocken{}, err
	}

	dtoTocken := dto.Tocken{
		Text: tocken.Text,
	}

	return dtoTocken, nil
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

func (a *app) SelectUserByTocken(dtoTocken dto.Tocken) (dto.User, error) {
	tocken := entity.Tocken{
		Text: dtoTocken.Text,
	}

	user, err := a.storage.SelectUserByTocken(tocken)
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

func (a *app) UpdateUser(dtoTocken dto.Tocken, dtoUser dto.User) error {
	tocken := entity.Tocken{
		Text: dtoTocken.Text,
	}

	userId, err := a.storage.SelectUserIdByTocken(tocken.Text)
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

func (a *app) CreateMessage(chatId int, dtoTocken dto.Tocken, dtoMessage dto.Message) error {
	//tocken := entity.Tocken{
	//	Text: dtoTocken.Text,
	//}
	//
	//user, err := a.storage.SelectUserByTocken(tocken)
	//if err != nil {
	//	return err
	//}

	//message := entity.Message{
	//	UserId:   user.Id,
	//	Type:     dtoMessage.Type,
	//	Text:     dtoMessage.Text,
	//	Created:  dtoMessage.Created,
	//	Messages: dtoMessage.Messages,
	//	Comments: dtoMessage.Messages,
	//}
	//a.storage.CreateMessage(chatId)
	return nil
}

func (a *app) GetMessage(tocken dto.Tocken, messageId, chatId int) (dto.Message, error) {
	return dto.Message{}, nil
}

func (a *app) GetPrevMessages(tocken dto.Tocken, chatId, offset, count int) ([]dto.Message, error) {
	return nil, nil
}

func (a *app) UpdateMessage(tocken dto.Tocken, chatId int, message dto.Message) error {
	return nil
}

func (a *app) SharedMessage(tocken dto.Tocken, srcChatId, dstChatId int, messages []dto.Message) error {
	return nil
}

func (a *app) CreateComment(tocken dto.Tocken, messageId int, message dto.Message) error {
	return nil
}

func (a *app) CreateChat(tocken dto.Tocken, chatName string, chatType int, adminsIds, userIds []int) error {
	return nil
}

func (a *app) SelectChatById(id int) (dto.Chat, error) {
	return dto.Chat{}, nil
}

func (a *app) GetPrevChats(tocken dto.Tocken, offset, count int) ([]dto.Chat, error) {
	return nil, nil
}

func (a *app) UpdateChat(tocken dto.Tocken, chat dto.Chat) error {
	return nil
}
