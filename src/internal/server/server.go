package server

import (
	"context"
	"fmt"

	"github.com/callmemars1/setka/src/bot/src/internal/server/handlers"
	"github.com/callmemars1/setka/src/bot/src/internal/services"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Register(g *echo.Echo)
}

func Run(ctx context.Context) error {
	serviceCollection, err := services.Build(ctx)
	if err != nil {
		return err
	}

	e := echo.New()

	for _, handler := range []Handler{
		&handlers.HealthCheck{},
		&handlers.CreateUser{
			UsersStorage:  serviceCollection.UsersStorage,
			OcservManager: serviceCollection.OcservManager,
			CertsManager:  serviceCollection.CertificatesManager,
		},
		&handlers.BanUser{
			UsersStorage: serviceCollection.UsersStorage,
		},
		&handlers.UnbanUser{
			UsersStorage: serviceCollection.UsersStorage,
		},
	} {
		handler.Register(e)
	}

	serverAddr := fmt.Sprintf("%s:%d", serviceCollection.Configuration.Host, serviceCollection.Configuration.Port)
	go func() {
		if err := e.Start(serverAddr); err != nil {
			fmt.Println("Server error:", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down server...")

	// Graceful shutdown
	return e.Shutdown(context.Background())
}
