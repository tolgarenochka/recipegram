package app

import (
	"context"
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
	topic := "count"
	groupID := "unique_group_id"

	cons, err := consumer.InitKafkaConsumer(broker, groupID)
	if err != nil {
		log.Fatal("Can't init Kafka consumer. Reason: ", err)
	}

	defer cons.Consumer.Close()

	// Подписка на топик
	err = cons.SubscribeTopic(topic)
	if err != nil {
		log.Printf("Failed to subscribe to topic: %v\n", err)
	}

	// Чтение сообщений
	for {
		msg, err := cons.Consumer.ReadMessage(-1)
		if err == nil {
			log.Printf("Received message: %s\n", msg.Value)
		} else {
			log.Printf("Error reading message: %v\n", err)
		}
	}

	ctx.Done()
}
