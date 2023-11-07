// handlers/friend.go
package handlers

import (
	"minimal_sns_app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FriendHandler struct {
	Repo repository.FriendRepository
}

func NewFriendHandler(repo repository.FriendRepository) *FriendHandler {
	return &FriendHandler{Repo: repo}
}

// RegisterRoutes registers the routes for friend operations
func (h *FriendHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/get_friend_list", h.GetFriendList)
	// 他のルートもここで登録...
}

// GetFriendList handles GET requests to retrieve a user's friend list
func (h *FriendHandler) GetFriendList(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	friends, err := h.Repo.GetFriends(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get friends")
	}

	return c.JSON(http.StatusOK, friends)
}
