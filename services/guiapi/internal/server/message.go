package server

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func (s *Server) sendMessageKafka(topic string, messageText []byte) error {
	// Отправка сообщения
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          messageText,
	}
	err := s.producer.Produce(message, nil)

	return err
}
