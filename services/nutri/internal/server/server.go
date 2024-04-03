package server

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
	"log"
)

type Server struct {
	eg       *errgroup.Group
	ctx      context.Context
	dbWizard dbWiz
	consumer sarama.Consumer
}

type dbWiz struct {
	dbWizard *sqlx.DB
}

func Init(dbW *sqlx.DB, consumer sarama.Consumer) *Server {
	s := Server{}
	s.dbWizard = dbWiz{
		dbWizard: dbW,
	}
	s.consumer = consumer

	return &s
}

func (s *Server) Run(ctx context.Context) error {
	s.eg, s.ctx = errgroup.WithContext(ctx)

	s.eg.Go(func() error {
		topic := "count"

		// Подписка на топик
		partitionConsumer, err := s.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Failed to subscribe to topic: %v\n", err)
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
					if err = s.CountNutri(msg.Value); err != nil {
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
	})

	return s.eg.Wait()
}
