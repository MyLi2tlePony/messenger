package dto

type User struct {
	PublicId string

	Login      string
	FirstName  string
	SecondName string
	Created    string
}

type Message struct {
	Id     int
	UserId int

	Type    int
	Changed bool
	Read    bool
	Text    string
	Created string

	Messages []Message
	Comments []Message
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
	Text string
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
	Tocken Tocken
	User   User
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
