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

func (h *PostHandler) CreatePost(c echo.Context) error {
	userIDValue := c.Get(middleware.UserIDKey)
	userID, ok := userIDValue.(int64)
	if !ok {
		return c.String(http.StatusUnauthorized, "unauthorized",)
	}
	title := c.FormValue("title")
	content := c.FormValue("content")
	image := c.FormValue("image")
	isPrivate := c.FormValue("is_private") == "true"
	var imagePtr *string
	if image != "" {
		imagePtr = &image
	}
	_, err := h.postService.CreatePost(c.Request().Context(), title, content, imagePtr, userID,isPrivate,)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error(),)
	}
	return c.Redirect(http.StatusSeeOther, "/",)
}

func (h *PostHandler) ShowCreatePost(c echo.Context) error {
	return c.Render(http.StatusOK, "create_post.html", nil,)
}