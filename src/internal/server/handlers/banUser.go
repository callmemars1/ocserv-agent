package handlers

import (
	"github.com/callmemars1/setka/src/bot/src/internal/users"
	"github.com/labstack/echo/v4"
)

type BanUser struct {
	UsersStorage *users.Storage
}

func (h *BanUser) Register(g *echo.Echo) {
	g.POST("/users/:username/ban", h.handle)
}

func (h *BanUser) handle(c echo.Context) error {
	username := c.Param("username")

	existingUser, err := h.UsersStorage.Get(username)
	if err != nil {
		return c.String(500, err.Error())
	}
	if existingUser == nil {
		return c.String(404, "user not found")
	}

	existingUser.IsBanned = true

	if err := h.UsersStorage.Save(existingUser); err != nil {
		return c.String(500, err.Error())
	}

	return c.JSON(200, "user banned")
}
