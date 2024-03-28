package server

import (
	"context"
	"github.com/IBM/sarama"
	"log"
)

type Cons struct {
	Consumer sarama.Consumer
}

func (c *Cons) Run(ctx context.Context) error {
	topic := "count"

	// Подписка на топик
	partitionConsumer, err := c.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to subscribe to topic: %v\n", err)
	}
	defer partitionConsumer.Close()

	// Чтение сообщений
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-partitionConsumer.Messages():
			switch msg.Topic {
			case "count":
				if err = CountNutri(msg.Value); err != nil {
					log.Printf("Error while nutri counting: %s", err.Error())
				}
				//maybe other topics
				//case "blabla":
			}
			log.Printf("Received message: %s\n", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("Error reading message: %v\n", err)
		}
	}
}
