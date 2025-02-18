package handlers

import (
	"github.com/callmemars1/setka/src/bot/src/internal/users"
	"github.com/labstack/echo/v4"
)

type UnbanUser struct {
	UsersStorage *users.Storage
}

func (h *UnbanUser) Register(g *echo.Echo) {
	g.POST("/users/:username/unban", h.handle)
}

func (h *UnbanUser) handle(c echo.Context) error {
	username := c.Param("username")

	existingUser, err := h.UsersStorage.Get(username)
	if err != nil {
		return c.String(500, err.Error())
	}
	if existingUser == nil {
		return c.String(404, "user not found")
	}

	existingUser.IsBanned = false

	if err := h.UsersStorage.Save(existingUser); err != nil {
		return c.String(500, err.Error())
	}

	return c.JSON(200, "user unbanned")
}
