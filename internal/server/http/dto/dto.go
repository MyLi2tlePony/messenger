package dto

type User struct {
	PublicId string `json:"public_id,omitempty"`

	Login      string `json:"login,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	SecondName string `json:"second_name,omitempty"`
	Created    string `json:"created,omitempty"`
}

type Message struct {
	Id     int `json:"id,omitempty"`
	UserId int `json:"user_id,omitempty"`

	Type    int    `json:"type,omitempty"`
	Changed bool   `json:"changed,omitempty"`
	Read    bool   `json:"read,omitempty"`
	Text    string `json:"text,omitempty"`
	Created string `json:"created,omitempty"`

	Messages []Message `json:"messages,omitempty"`
	Comments []Message `json:"comments,omitempty"`
}

type Chat struct {
	ChatId string

	Type    int
	Created string
	Name    string

	AdminIds []int
	UserIds  []int
}

type Tocken struct {
	Text string `json:"text,omitempty"`
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
	Tocken Tocken `json:"tocken"`
	User   User   `json:"user"`
}

type SendMessageRequest struct {
	Tocken  Tocken
	ChatId  int
	Message Message
}

type ChangeMessageRequest struct {
	Tocken  Tocken
	ChatId  int
	Message Message
}

type GetMessageRequest struct {
	Tocken    Tocken
	ChatId    int
	MessageId int
}

type GetPrevMessagesRequest struct {
	Tocken Tocken
	ChatId int
	Offset int
	Count  int
}

type UpdateMessageRequest struct {
	Tocken  Tocken
	ChatId  int
	Message Message
}

type SharedMessageRequest struct {
	Tocken    Tocken
	SrcChatId int
	DstChatId int
	Messages  []Message
}

type CreateCommentRequest struct {
	Tocken    Tocken
	MessageId int
	Message   Message
}

type CreateChatRequest struct {
	Tocken Tocken
	Name   string
	Type   int

	UserIds  []int
	AdminIds []int
}

type GetPrevChatsRequest struct {
	Tocken Tocken
	Offset int
	Count  int
}

type UpdateChatRequest struct {
	Tocken Tocken
	Chat   Chat
}
