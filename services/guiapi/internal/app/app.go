package app

import (
	"context"
	"guiapi/internal/dbwizard"
	"guiapi/internal/server"
	"log"
	"os/signal"
	"syscall"
)

func Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	log.Println("Init service guiapi...")

	s := server.NewServer()
	s.Init()

	dbWizard, err := dbwizard.NewConnect()
	if err != nil {
		log.Fatal("Database wizard init failed. Reason:", err)
	}
	defer func() {
		if err = dbWizard.Quit(); err != nil {
			log.Fatal("Error while db connecting:", err)
		}
	}()

	err = s.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
