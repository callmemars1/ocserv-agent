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
		return c.JSON(500, echo.Map{"message": err.Error()})
	}
	if existingUser == nil {
		return c.JSON(404, echo.Map{"message": "user not found"})
	}

	existingUser.IsBanned = true

	if err := h.UsersStorage.Save(existingUser); err != nil {
		return c.JSON(500, echo.Map{"message": err.Error()})
	}

	return c.JSON(200, echo.Map{"message": "user banned", "username": username})
}
