package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"weblog/internal/services"
)

type AuthHandler struct {
	userService    *service.UserService
	sessionService *service.SessionService
}

func NewAuthHandler(userService *service.UserService, sessionService *service.SessionService,) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		sessionService: sessionService,
	}
}

func (h *AuthHandler) Signup(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	user, err := h.userService.Register(c.Request().Context(), username, password,)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	token, err := h.sessionService.CreateSession(c.Request().Context(), user.ID,)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create session",)
	}
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = 24 * 60 * 60
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	user, err := h.userService.Login(c.Request().Context(), username, password,)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}
	token, err := h.sessionService.CreateSession(c.Request().Context(), user.ID,)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create session",)
	}
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = 24 * 60 * 60
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthHandler) Logout(c echo.Context) error {
	cookie, err := c.Cookie("session_token")
	if err == nil {
		err = h.sessionService.DeleteSession(c.Request().Context(), cookie.Value,)
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to logout",)
		}
	}
	clearCookie := new(http.Cookie)
	clearCookie.Name = "session_token"
	clearCookie.Value = ""
	clearCookie.Path = "/"
	clearCookie.HttpOnly = true
	clearCookie.MaxAge = -1
	c.SetCookie(clearCookie)
	return c.Redirect(http.StatusSeeOther, "/login")
}