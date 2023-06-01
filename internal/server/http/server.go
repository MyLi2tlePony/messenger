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
	SelectUserByPublicId(publicId string) (dto.User, error)
	SelectUserByTocken(dtoTocken dto.Tocken) (dto.User, error)
	UpdateUser(tocken dto.Tocken, user dto.User) error

	CreateTocken(request dto.CreateTockenRequest) (dto.Tocken, error)

	CreateMessage(chatPublicId int, tocken dto.Tocken, message dto.Message) error
	GetMessage(tocken dto.Tocken, messageId, chatId int) (dto.Message, error)
	GetPrevMessages(tocken dto.Tocken, chatId, offset, count int) ([]dto.Message, error)
	UpdateMessage(tocken dto.Tocken, chatId int, message dto.Message) error
	SharedMessage(tocken dto.Tocken, srcChatId, dstChatId int, messages []dto.Message) error
	CreateComment(tocken dto.Tocken, messageId int, message dto.Message) error

	CreateChat(tocken dto.Tocken, chatName string, chatType int, adminsIds, userIds []int) error
	SelectChatById(id int) (dto.Chat, error)
	GetPrevChats(tocken dto.Tocken, offset, count int) ([]dto.Chat, error)
	UpdateChat(tocken dto.Tocken, chat dto.Chat) error
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

func (s *Server) CreateTocken(ctx echo.Context) error {
	body := new(dto.CreateTockenRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	tocken, err := s.app.CreateTocken(*body)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, tocken)
}

func (s *Server) SelectUserByTocken(ctx echo.Context) error {
	body := new(dto.Tocken)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	user, err := s.app.SelectUserByTocken(*body)
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

	err := s.app.UpdateUser(body.Tocken, body.User)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) CreateMessage(ctx echo.Context) error {
	body := new(dto.SendMessageRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.CreateMessage(body.ChatId, body.Tocken, body.Message)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) GetMessage(ctx echo.Context) error {
	body := new(dto.GetMessageRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	message, err := s.app.GetMessage(body.Tocken, body.MessageId, body.ChatId)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, message)
}

func (s *Server) GetPrevMessages(ctx echo.Context) error {
	body := new(dto.GetPrevMessagesRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	messages, err := s.app.GetPrevMessages(body.Tocken, body.ChatId, body.Offset, body.Count)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, messages)
}

func (s *Server) UpdateMessage(ctx echo.Context) error {
	body := new(dto.UpdateMessageRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.UpdateMessage(body.Tocken, body.ChatId, body.Message)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) SharedMessage(ctx echo.Context) error {
	body := new(dto.SharedMessageRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.SharedMessage(body.Tocken, body.SrcChatId, body.DstChatId, body.Messages)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) CreateComment(ctx echo.Context) error {
	body := new(dto.CreateCommentRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.CreateComment(body.Tocken, body.MessageId, body.Message)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) CreateChat(ctx echo.Context) error {
	body := new(dto.CreateChatRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.CreateChat(body.Tocken, body.Name, body.Type, body.AdminIds, body.UserIds)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) GetChat(ctx echo.Context) error {
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

func (s *Server) GetPrevChats(ctx echo.Context) error {
	body := new(dto.GetPrevChatsRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	chats, err := s.app.GetPrevChats(body.Tocken, body.Offset, body.Count)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, chats)
}

func (s *Server) UpdateChat(ctx echo.Context) error {
	body := new(dto.UpdateChatRequest)
	if err := ctx.Bind(body); err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := s.app.UpdateChat(body.Tocken, body.Chat)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusOK)
}
