package tasks

import (
	"context"
	"time"

	"github.com/callmemars1/setka/src/bot/src/internal/ocserv"
	"github.com/callmemars1/setka/src/bot/src/internal/users"
)

type Ban struct {
	UsersStorage  *users.Storage
	OcservManager *ocserv.Manager
}

func (t *Ban) Run(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			users, err := t.UsersStorage.GetAll()
			if err != nil {
				continue
			}
			for _, user := range users {
				if user.IsBanned {
					if err := t.OcservManager.DisconnectUser(user.Username); err != nil {
						continue
					}
				}
			}
		case <-ctx.Done():
			return ctx.Err() // Graceful shutdown
		}
	}
}
