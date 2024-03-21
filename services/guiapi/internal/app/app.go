package app

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"guiapi/internal/server"
	"log"
	"os/signal"
	"syscall"
)

func Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	log.Println("Init service guiapi...")

	dbWizard, err := sqlx.ConnectContext(context.Background(), "postgres", "postgres://postgres:1q2w3e4r5t@postgres:5432/recipegram?sslmode=disable")
	if err != nil {
		log.Fatal("Database wizard init failed. Reason:", err)
	}

	defer func() {
		if err = dbWizard.Close(); err != nil {
			log.Fatal("Error close db connecting:", err)
		}
	}()

	broker := "kafka:9093"
	//topic := "your_topic_name"

	// Настройка продюсера
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatal("Kafka producer init failed. Reason:", err)
	}
	defer producer.Close()

	s := server.NewServer()
	s.Init(dbWizard, producer)

	err = s.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
