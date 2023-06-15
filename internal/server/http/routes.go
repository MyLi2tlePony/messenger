package httpsrv

import "github.com/MyLi2tlePony/messenger/internal/server/http/urls"

func SetupRoutes(s *Server) {
	s.echo.GET(urls.UrlPing, s.ping)

	s.echo.GET(urls.UrlUserId, s.SelectUserByPublicId)
	s.echo.POST(urls.UrlGetUser, s.SelectUserByToken)
	s.echo.POST(urls.UrlUser, s.CreateUser)
	s.echo.POST(urls.UrlUpdateUser, s.UpdateUser)
	s.echo.POST(urls.UrlToken, s.CreateToken)

	//
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//

	s.echo.POST(urls.UrlCreateChat, s.CreateChat)
	s.echo.GET(urls.UrlSelectChatById, s.SelectChatById)

	s.echo.POST(urls.UrlCreateMessage, s.CreateMessage)
	s.echo.POST(urls.UrlSelectTopMessages, s.SelectTopMessages)
	s.echo.POST(urls.UrlSelectMessagesByIds, s.SelectMessagesById)
	s.echo.POST(urls.UrlDeleteMessage, s.DeleteMessage)

	s.echo.POST(urls.UrlGetUserChats, s.GetUserChats)
	s.echo.POST(urls.UrlDeleteUserChat, s.DeleteUserChat)

	s.echo.POST(urls.UrlCreateParticipant, s.CreateParticipant)
	s.echo.GET(urls.UrlSelectParticipant, s.SelectParticipant)
	s.echo.POST(urls.UrlDeleteParticipant, s.DeleteParticipant)
}
