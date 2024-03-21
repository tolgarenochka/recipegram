package app

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
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
	broker := "kafka:9093"
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
		fmt.Printf("Failed to subscribe to topic: %v\n", err)
	}

	// Чтение сообщений
	for {
		msg, err := cons.Consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message: %s\n", msg.Value)
		} else {
			fmt.Printf("Error reading message: %v\n", err)
		}
	}

	ctx.Done()
}
