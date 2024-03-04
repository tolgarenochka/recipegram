package app

import (
	"context"
	"github.com/tolgarenochka/recipegram/db/dbwizard"
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

	dbWizard, err := dbwizard.NewConnect()
	if err != nil {
		log.Fatal("Database wizard init failed. Reason:", err)
	}
	defer func() {
		if err = dbWizard.Quit(); err != nil {
			logger.Fatal(err.Error())
		}
	}()

	err := s.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
