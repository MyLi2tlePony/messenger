package httpsrv

import "github.com/MyLi2tlePony/messenger/internal/server/http/urls"

func SetupRoutes(s *Server) {
	s.echo.GET(urls.UrlPing, s.ping)

	s.echo.GET(urls.UrlUserId, s.SelectUserByPublicId)
	s.echo.GET(urls.UrlUser, s.SelectUserByTocken)
	s.echo.POST(urls.UrlUser, s.CreateUser)
	s.echo.PATCH(urls.UrlUser, s.UpdateUser)

	s.echo.POST(urls.UrlTocken, s.CreateTocken)
}
