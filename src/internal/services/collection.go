package services

import (
	"context"
	"log/slog"

	"github.com/callmemars1/setka/src/bot/src/internal/certs"
	"github.com/callmemars1/setka/src/bot/src/internal/configuration"
	"github.com/callmemars1/setka/src/bot/src/internal/logger"
	"github.com/callmemars1/setka/src/bot/src/internal/ocserv"
	"github.com/callmemars1/setka/src/bot/src/internal/users"
)

type Collection struct {
	Logger              *slog.Logger
	Configuration       *configuration.C
	OcservManager       *ocserv.Manager
	CertificatesManager *certs.Manager
	UsersStorage        *users.Storage
}

func Build(ctx context.Context) (*Collection, error) {
	collection := &Collection{}

	configuration, err := configuration.Initialize()
	if err != nil {
		return nil, err
	}
	collection.Configuration = &configuration

	logger, err := logger.BuildForSyslog()
	if err != nil {
		return nil, err
	}
	collection.Logger = logger

	usersStorage, err := users.NewStorage(configuration.DataDirectory)
	if err != nil {
		return nil, err
	}
	collection.UsersStorage = usersStorage

	collection.OcservManager = ocserv.NewManager(configuration.OcservPasswdPath)
	collection.CertificatesManager = certs.NewManager(certs.Configuration{
		Organization:         configuration.Organization,
		DataDirectory:        configuration.DataDirectory,
		CaCertificatePath:    configuration.CaCertificatePath,
		CaPrivateKeyPath:     configuration.CaPrivateKeyPath,
		ClientPrivateKeyPath: configuration.ClientPrivateKeyPath,
	})

	return collection, nil
}
