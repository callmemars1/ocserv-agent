package handlers

import (
	"encoding/base64"

	"github.com/callmemars1/setka/src/bot/src/internal/certs"
	"github.com/callmemars1/setka/src/bot/src/internal/ocserv"
	"github.com/callmemars1/setka/src/bot/src/internal/users"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type GetUserResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsBanned bool   `json:"is_banned"`

	CertificateBase64 string `json:"certificate_base64"`
}

type GetUser struct {
	UsersStorage  *users.Storage
	OcservManager *ocserv.Manager
	CertsManager  *certs.Manager
}

func (h *GetUser) Register(g *echo.Echo) {
	g.GET("/users/:username", h.handle)
}

func (h *GetUser) handle(c echo.Context) error {
	username := c.Param("username")

	existingUser, err := h.UsersStorage.Get(username)
	if err != nil {
		return c.String(500, err.Error())
	}

	if existingUser == nil {
		return c.JSON(404, echo.Map{
			"message": "user not found",
		})
	}

	c.Logger().Infoj(log.JSON{"message": "user exists", "username": username})
	certificateBytes, err := h.CertsManager.ReadClientCertificateP12(existingUser.Username)
	if err != nil {
		return c.String(500, err.Error())
	}

	return c.JSON(200, CreateUserResponse{
		Created:  false,
		Username: existingUser.Username,
		Password: existingUser.Password,
		IsBanned: existingUser.IsBanned,

		CertificateBase64: base64.StdEncoding.EncodeToString(certificateBytes),
	})
}
