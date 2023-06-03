package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/MyLi2tlePony/messenger/internal/server/http/dto"
	"github.com/MyLi2tlePony/messenger/internal/server/http/urls"
	"io"
	"net/http"
)

var (
	ErrBadRequest = errors.New("error bad request")
)

type client struct {
	httpClient http.Client
	domain     string
}

func NewHttpClient(httpClient http.Client, domain string) *client {
	return &client{
		httpClient: httpClient,
		domain:     domain,
	}
}

func (c *client) CreateUser(createUserRequest dto.CreateUserRequest) error {
	content, err := json.Marshal(createUserRequest)
	if err != nil {
		return err
	}

	url := c.domain + urls.UrlUser

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(content))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrBadRequest
	}

	return nil
}

func (c *client) CreateTocken(createTockenRequest dto.CreateTockenRequest) (dto.Tocken, error) {
	content, err := json.Marshal(createTockenRequest)
	if err != nil {
		return dto.Tocken{}, err
	}

	url := c.domain + urls.UrlTocken

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(content))
	if err != nil {
		return dto.Tocken{}, err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return dto.Tocken{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return dto.Tocken{}, ErrBadRequest
	}

	content, err = io.ReadAll(resp.Body)
	if err != nil {
		return dto.Tocken{}, err
	}

	err = resp.Body.Close()
	if err != nil {
		return dto.Tocken{}, err
	}

	tocken := dto.Tocken{}
	err = json.Unmarshal(content, &tocken)
	if err != nil {
		return dto.Tocken{}, err
	}

	return tocken, nil
}

func (c *client) SelectUserByPublicId(userId string) (dto.User, error) {
	url := c.domain + urls.UrlUser + "/" + userId
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return dto.User{}, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return dto.User{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return dto.User{}, ErrBadRequest
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.User{}, err
	}

	err = resp.Body.Close()
	if err != nil {
		return dto.User{}, err
	}

	user := dto.User{}
	err = json.Unmarshal(content, &user)
	if err != nil {
		return dto.User{}, err
	}

	return user, nil
}

func (c *client) SelectUserByTocken(dtoTocken dto.Tocken) (dto.User, error) {
	content, err := json.Marshal(dtoTocken)
	if err != nil {
		return dto.User{}, err
	}

	url := c.domain + urls.UrlGetUser
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(content))
	if err != nil {
		return dto.User{}, err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return dto.User{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return dto.User{}, ErrBadRequest
	}

	content, err = io.ReadAll(resp.Body)
	if err != nil {
		return dto.User{}, err
	}

	err = resp.Body.Close()
	if err != nil {
		return dto.User{}, err
	}

	user := dto.User{}
	err = json.Unmarshal(content, &user)
	if err != nil {
		return dto.User{}, err
	}

	return user, nil
}

func (c *client) UpdateUser(updateUserRequest dto.UpdateUserRequest) error {
	content, err := json.Marshal(updateUserRequest)
	if err != nil {
		return err
	}

	url := c.domain + urls.UrlUser
	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(content))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrBadRequest
	}

	return nil
}

//func (c *client) CreateMessage(ctx echo.Context) error {
//}
//
//func (c *client) GetMessage(ctx echo.Context) error {
//}
//
//func (c *client) GetPrevMessages(ctx echo.Context) error {
//}
//
//func (c *client) UpdateMessage(ctx echo.Context) error {
//
//}
//
//func (c *client) SharedMessage(ctx echo.Context) error {
//}
//
//func (c *client) CreateComment(ctx echo.Context) error {
//}
//
//func (c *client) CreateChat(ctx echo.Context) error {
//}
//
//func (c *client) GetChat(ctx echo.Context) error {
//}
//
//func (c *client) GetPrevChats(ctx echo.Context) error {
//}
//
//func (c *client) UpdateChat(ctx echo.Context) error {
//}
