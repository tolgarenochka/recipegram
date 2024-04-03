package app

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"nutri/internal/server"
	"os/signal"
	"syscall"
)

func Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	log.Println("Init service nutri...")

	dbWizard, err := sqlx.ConnectContext(context.Background(), "postgres", "postgres://postgres:1q2w3e4r5t@postgres:5432/recipegram?sslmode=disable")
	if err != nil {
		log.Fatal("Database wizard init failed. Reason:", err)
	}

	defer func() {
		if err = dbWizard.Close(); err != nil {
			log.Fatal("Error close db connecting:", err)
		}
	}()

	// TODO: здесь и в других местах вынести такие параметры в конфиг файл
	// Адрес и порт брокера Kafka
	broker := "kafka:9092"
	//groupID := "group_id"

	// Настройка консьюмера
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	cons, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log.Fatal("Can't init Kafka consumer. Reason: ", err)
	}
	defer cons.Close()

	s := server.Init(dbWizard, cons)
	err = s.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
