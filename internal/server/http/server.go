package httpsrv

import (
	"errors"
	"github.com/MyLi2tlePony/messenger/internal/server/http/dto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
)

type App interface {
	CreateUser(request dto.CreateUserRequest) error
	CreateToken(request dto.CreateTokenRequest) (dto.Token, error)
	SelectUserByPublicId(publicId string) (dto.User, error)
	SelectUserByToken(dtoToken dto.Token) (dto.User, error)
	UpdateUser(dtoToken dto.Token, dtoUser dto.User) error

	//
	////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//

	CreateChat(createChatRequest dto.CreateChatRequest) error
	SelectChatById(id int) (dto.Chat, error)
	CreateMessage(chatId int, dtoToken dto.Token, dtoMessage dto.Message) error
	SelectTopMessages(dtoToken dto.Token, chatId, limit int) ([]dto.Message, error)
	SelectMessagesById(chatId, minId, maxId int, dtoToken dto.Token) ([]dto.Message, error)
	DeleteMessage(chatId, messageId int, dtoToken dto.Token) error
	GetUserChats(dtoToken dto.Token) ([]dto.Chat, error)
	DeleteUserChat(dtoToken dto.Token, chatId int) error
	CreateParticipant(dtoToken dto.Token, chatId int, dtoParticipant dto.Participant) error
	SelectParticipant(chatId int) ([]dto.Participant, error)
	DeleteParticipant(dtoToken dto.Token, chatId, participantId int) error
}

var (
	ErrNoRows = errors.New("no rows in result set")
)

type Server struct {
	app App

	echo *echo.Echo
}

func NewServer(lvl log.Lvl, app App) *Server {
	server := &Server{
		app:  app,
		echo: echo.New(),
	}

	server.echo.Logger.SetLevel(lvl)
	SetupRoutes(server)

	return server
}

func (s *Server) Start(address string) error {
	err := s.echo.Start(address)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	err := s.echo.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) ping(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "pong")
}

func (s *Server) CreateUser(ctx echo.Context) error {
	body := new(dto.CreateUserRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.CreateUser(*body)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) CreateToken(ctx echo.Context) error {
	body := new(dto.CreateTokenRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	token, err := s.app.CreateToken(*body)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, token)
}

func (s *Server) SelectUserByToken(ctx echo.Context) error {
	body := new(dto.Token)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	user, err := s.app.SelectUserByToken(*body)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, user)
}

func (s *Server) SelectUserByPublicId(ctx echo.Context) error {
	id := ctx.Param("id")

	user, err := s.app.SelectUserByPublicId(id)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, user)
}

func (s *Server) UpdateUser(ctx echo.Context) error {
	body := new(dto.UpdateUserRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.UpdateUser(body.Token, body.User)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

//
//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//

func (s *Server) CreateChat(ctx echo.Context) error {
	body := new(dto.CreateChatRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.CreateChat(dto.CreateChatRequest{})
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) SelectChatById(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	user, err := s.app.SelectChatById(id)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, user)
}

func (s *Server) CreateMessage(ctx echo.Context) error {
	body := new(dto.CreateMessageRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.CreateMessage(body.ChatId, body.Token, body.Message)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) SelectTopMessages(ctx echo.Context) error {
	body := new(dto.SelectTopMessagesRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	messages, err := s.app.SelectTopMessages(body.Token, body.ChatId, body.Limit)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, messages)
}

func (s *Server) SelectMessagesById(ctx echo.Context) error {
	body := new(dto.SelectMessagesByIdRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	messages, err := s.app.SelectMessagesById(body.ChatId, body.MinId, body.MaxId, body.Token)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, messages)
}

func (s *Server) DeleteMessage(ctx echo.Context) error {
	body := new(dto.DeleteMessageRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.DeleteMessage(body.ChatId, body.MessageId, body.Token)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) GetUserChats(ctx echo.Context) error {
	body := new(dto.GetUserChatsRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	result, err := s.app.GetUserChats(body.Token)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, result)
}

func (s *Server) DeleteUserChat(ctx echo.Context) error {
	body := new(dto.DeleteUserChatRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.DeleteUserChat(body.Token, body.ChatId)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) CreateParticipant(ctx echo.Context) error {
	body := new(dto.CreateParticipantRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.CreateParticipant(body.Token, body.ChatId, body.Participant)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) SelectParticipant(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("chat_id"))
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	results, err := s.app.SelectParticipant(id)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusBadRequest, results)
	}

	return ctx.JSON(http.StatusOK, results)
}

func (s *Server) DeleteParticipant(ctx echo.Context) error {
	body := new(dto.DeleteParticipantRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.DeleteParticipant(body.Token, body.ChatId, body.ParticipantId)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}
