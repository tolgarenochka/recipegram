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
	//// Отправка сообщения
	//message := &kafka.Message{
	//	TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	//	Value:          messageText,
	//}
	//err := s.producer.Produce(message, nil)

	return err
}
