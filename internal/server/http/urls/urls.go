package urls

const (
	UrlPing = "/ping"

	UrlToken   = "/token"
	UrlGetUser = "/get/user"
	UrlUser    = "/user"
	UrlUserId  = "/user/:id"

	UrlCreateChat     = "/chat"
	UrlSelectChatById = "/chat/:id"

	UrlCreateMessage       = "/message"
	UrlSelectTopMessages   = "/messages/top"
	UrlSelectMessagesByIds = "/messages/ids"
	UrlDeleteMessage       = "/messages/delete"

	UrlGetUserChats   = "/user/chat"
	UrlDeleteUserChat = "/user/chat/delete"

	UrlCreateParticipant = "/participant"
	UrlSelectParticipant = "/participant/:chat_id"
	UrlDeleteParticipant = "/participant/delete"
)
