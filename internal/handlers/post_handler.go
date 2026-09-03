package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"weblog/internal/middleware"
	"weblog/internal/services"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService,) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) Feed(c echo.Context) error {
	userIDValue := c.Get(middleware.UserIDKey)
	userID, ok := userIDValue.(int64)
	if !ok {
		return c.String(http.StatusUnauthorized, "unauthorized",)
	}

	posts, err := h.postService.GetFeed(c.Request().Context(), userID,)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to load posts",)
	}

	return c.Render(http.StatusOK, "feed.html",map[string]interface{} {
			"Posts": posts,
		},)
}