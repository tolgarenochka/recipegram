package app

import (
	"context"
	"guiapi/internal/server"
	"log"
	"os/signal"
	"syscall"
)

func Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s := server.NewServer()
	s.Init()

	err := s.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
