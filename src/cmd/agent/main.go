package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/callmemars1/setka/src/bot/src/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx); err != nil {
		panic(err)
	}
}
