package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"weblog/internal/services"
)

const UserIDKey = "user_id"
const SessionCookieName = "session_token"

type AuthMiddleware struct {
	sessionService *service.SessionService
}

func NewAuthMiddleware(sessionService *service.SessionService,) *AuthMiddleware {
	return &AuthMiddleware{
		sessionService: sessionService,
	}
}

func (m *AuthMiddleware) RequireAuth(next echo.HandlerFunc,) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(SessionCookieName)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		session, err := m.sessionService.GetSession(context.Background(), cookie.Value,)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		c.Set(UserIDKey, session.UserID)
		return next(c)
	}
}