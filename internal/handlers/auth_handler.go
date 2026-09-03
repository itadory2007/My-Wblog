package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"weblog/internal/service"
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