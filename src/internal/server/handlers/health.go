package handlers

import (
	"github.com/callmemars1/setka/src/bot/src/internal/users"
	"github.com/labstack/echo/v4"
)

type HealthCheck struct {
	UsersStorage *users.Storage
}

func (h *HealthCheck) Register(g *echo.Echo) {
	g.GET("/checks/health", h.handle)
}

func (h *HealthCheck) handle(c echo.Context) error {
	return c.JSON(200, echo.Map{"status": "ok"})
}
