// handlers/friend.go
package handlers

import (
	"minimal_sns_app/logutils"
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
	e.GET("/get_friend_of_friend_list", h.Get2hopsFriends)
	e.DELETE("/delete_friend", h.DeleteFriend)
}

// GetFriendList handles GET requests to retrieve a user's friend list
func (h *FriendHandler) GetFriendList(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	friends, err := h.Repo.GetFriends(userID)
	if err != nil {
		logutils.Error("Failed to get friends")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get friends")
	}

	return c.JSON(http.StatusOK, friends)
}

// Get2hopsFriends handles GET requests to retrieve a user's friend list
func (h *FriendHandler) Get2hopsFriends(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	friends, err := h.Repo.Get2hopsFriends(userID)
	if err != nil {
		logutils.Error("Failed to get friends")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get friends")
	}

	return c.JSON(http.StatusOK, friends)
}

// DeleteFriend handles DELETE requests to delete a friend
func (h *FriendHandler) DeleteFriend(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	friendIDParam := c.QueryParam("friend_id")
	friendID, err := strconv.ParseInt(friendIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid friend_id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid friend_id")
	}

	err = h.Repo.DeleteFriend(userID, friendID)
	if err != nil {
		logutils.Error("Failed to delete friend")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete friend")
	}

	// successを返す
	return c.String(http.StatusOK, "success")
}
