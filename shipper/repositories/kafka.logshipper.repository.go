package repositories

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaLogShipperRepository struct {
	producer *kafka.Producer
}

func NewKafkaLogShipperRepository(producer *kafka.Producer) *KafkaLogShipperRepository {
	return &KafkaLogShipperRepository{producer: producer}
}

func (k KafkaLogShipperRepository) Send(topic string, value []byte) error {
	// Produce messages to topic (asynchronously)
	err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)

	k.producer.Flush(10)

	return err
}