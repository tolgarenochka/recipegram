package server

import (
	"github.com/IBM/sarama"
)

func (s *Server) sendMessageKafka(topic string, messageText string) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(messageText),
	}

	// Отправка сообщения
	_, _, err := s.producer.SendMessage(message)

	return err
}
