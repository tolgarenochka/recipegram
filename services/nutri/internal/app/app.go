package app

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"nutri/internal/consumer"
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

	// Адрес и порт брокера Kafka
	broker := "kafka:9092"
	groupID := "group_id"

	cons, err := consumer.InitKafkaConsumer(broker, groupID)
	if err != nil {
		log.Fatal("Can't init Kafka consumer. Reason: ", err)
	}

	defer cons.Consumer.Close()

	topic := "count"

	// Подписка на топик
	err = cons.SubscribeTopic(topic)
	if err != nil {
		log.Printf("Failed to subscribe to topic: %v\n", err)
	}
	// Чтение сообщений
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := cons.Consumer.ReadMessage(100)
			if err == nil {
				log.Printf("Received message: %s\n", msg.String())
			} else if !err.(kafka.Error).IsTimeout() {
				log.Printf("Error reading message: %v\n", err)
			}
		}
	}

}
