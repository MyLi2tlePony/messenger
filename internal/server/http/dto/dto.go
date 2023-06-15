package dto

import "time"

type User struct {
	PublicId string `json:"public_id,omitempty"`

	Login      string    `json:"login,omitempty"`
	FirstName  string    `json:"first_name,omitempty"`
	SecondName string    `json:"second_name,omitempty"`
	Created    time.Time `json:"created,omitempty"`
}

type Message struct {
	Id     int `json:"id,omitempty"`
	UserId int `json:"user_id,omitempty"`

	Changed bool `json:"changed,omitempty"`

	Text               string `json:"text,omitempty"`
	CommentedMessageId int    `json:"commented_message_id,omitempty"`

	Created time.Time `json:"created"`
}

type Chat struct {
	Id int `json:"id,omitempty"`

	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`

	Open    bool      `json:"open,omitempty"`
	Created time.Time `json:"created"`
}

type Participant struct {
	Id     int `json:"id,omitempty"`
	UserId int `json:"user_id,omitempty"`

	Write   bool `json:"write,omitempty"`
	Post    bool `json:"post,omitempty"`
	Comment bool `json:"comment,omitempty"`
	Delete  bool `json:"delete,omitempty"`

	AddParticipant    bool `json:"add_participant,omitempty"`
	DeleteParticipant bool `json:"delete_participant,omitempty"`
}

type Token struct {
	Text string `json:"text,omitempty"`
}

type CreateChatRequest struct {
	Chat         Chat          `json:"chat"`
	Participants []Participant `json:"participants,omitempty"`
}

type CreateUserRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type CreateTockenRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type UpdateUserRequest struct {
	Token Token `json:"tocken"`
	User  User  `json:"user"`
}

type CreateMessageRequest struct {
	Token   Token   `json:"token"`
	ChatId  int     `json:"chat_id,omitempty"`
	Message Message `json:"message"`
}

type ChangeMessageRequest struct {
	Token   Token   `json:"token"`
	ChatId  int     `json:"chat_id,omitempty"`
	Message Message `json:"message"`
}

type GetMessageRequest struct {
	Token     Token `json:"token"`
	ChatId    int   `json:"chat_id,omitempty"`
	MessageId int   `json:"message_id,omitempty"`
}

type GetPrevMessagesRequest struct {
	Token  Token `json:"token"`
	ChatId int   `json:"chat_id,omitempty"`
	Offset int   `json:"offset,omitempty"`
	Count  int   `json:"count,omitempty"`
}

type UpdateMessageRequest struct {
	Token   Token   `json:"token"`
	ChatId  int     `json:"chat_id,omitempty"`
	Message Message `json:"message"`
}

type SharedMessageRequest struct {
	Token     Token     `json:"token"`
	SrcChatId int       `json:"src_chat_id,omitempty"`
	DstChatId int       `json:"dst_chat_id,omitempty"`
	Messages  []Message `json:"messages,omitempty"`
}

type CreateCommentRequest struct {
	Token     Token   `json:"token"`
	MessageId int     `json:"message_id,omitempty"`
	Message   Message `json:"message"`
}

type GetPrevChatsRequest struct {
	Token  Token `json:"token"`
	Offset int   `json:"offset,omitempty"`
	Count  int   `json:"count,omitempty"`
}

type UpdateChatRequest struct {
	Token Token `json:"token"`
	Chat  Chat  `json:"chat"`
}

type SelectTopMessagesRequest struct {
	Token  Token `json:"token"`
	ChatId int   `json:"chat_id,omitempty"`
	Limit  int   `json:"limit,omitempty"`
}

type SelectMessagesByIdRequest struct {
	ChatId int   `json:"chat_id,omitempty"`
	MinId  int   `json:"min_id,omitempty"`
	MaxId  int   `json:"max_id,omitempty"`
	Token  Token `json:"token"`
}

type DeleteMessageRequest struct {
	ChatId    int   `json:"chat_id,omitempty"`
	MessageId int   `json:"message_id,omitempty"`
	Token     Token `json:"token"`
}

type GetUserChatsRequest struct {
	Token Token `json:"token"`
}

type DeleteUserChatRequest struct {
	Token  Token `json:"token"`
	ChatId int   `json:"chat_id,omitempty"`
}

type CreateParticipantRequest struct {
	Token       Token       `json:"token"`
	ChatId      int         `json:"chat_id,omitempty"`
	Participant Participant `json:"participant"`
}

type DeleteParticipantRequest struct {
	Token         Token `json:"token"`
	ChatId        int   `json:"chat_id,omitempty"`
	ParticipantId int   `json:"participant_id,omitempty"`
}
