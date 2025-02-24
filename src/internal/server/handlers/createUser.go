package handlers

import (
	"encoding/base64"
	"math/rand"

	"github.com/callmemars1/setka/src/bot/src/internal/certs"
	"github.com/callmemars1/setka/src/bot/src/internal/ocserv"
	"github.com/callmemars1/setka/src/bot/src/internal/users"
	"github.com/labstack/echo/v4"
)

type CreateUserResponse struct {
	Created bool `json:"created"`

	Username string `json:"username"`
	Password string `json:"password"`
	IsBanned bool   `json:"is_banned"`

	CertificateBase64 string `json:"certificate_base64"`
}

type CreateUser struct {
	UsersStorage  *users.Storage
	OcservManager *ocserv.Manager
	CertsManager  *certs.Manager
}

func (h *CreateUser) Register(g *echo.Echo) {
	g.PUT("/users/:username", h.handle)
}

func (h *CreateUser) handle(c echo.Context) error {
	username := c.Param("username")

	existingUser, err := h.UsersStorage.Get(username)
	if err != nil {
		return c.String(500, err.Error())
	}

	if existingUser != nil {
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

	user := &users.User{
		Username: username,
		Password: generateNumericPassword(6),
		IsBanned: false,
	}

	if err := h.OcservManager.AddUser(user.Username, user.Password); err != nil {
		return c.String(500, err.Error())
	}

	certificateBytes, err := h.CertsManager.IssueCertificateForUser(user)
	if err != nil {
		return c.String(500, err.Error())
	}

	if err := h.UsersStorage.Save(user); err != nil {
		return c.String(500, err.Error())
	}

	return c.JSON(200, CreateUserResponse{
		Created:  true,
		Username: user.Username,
		Password: user.Password,
		IsBanned: user.IsBanned,

		CertificateBase64: base64.StdEncoding.EncodeToString(certificateBytes),
	})
}

func generateNumericPassword(length int) string {
	digits := "0123456789"
	password := make([]byte, length)
	for i := range password {
		password[i] = digits[rand.Intn(len(digits))]
	}
	return string(password)
}
