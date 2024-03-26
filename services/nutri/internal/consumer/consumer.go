package consumer

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type Cons struct {
	Consumer *kafka.Consumer
}

func InitKafkaConsumer(broker, group string) (*Cons, error) {
	// Настройка консьюмера
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          group,
		"auto.offset.reset": "smallest",
	})

	return &Cons{consumer}, err
}

func (c *Cons) SubscribeTopic(topic string) error {
	return c.Consumer.SubscribeTopics([]string{topic}, nil)
}
