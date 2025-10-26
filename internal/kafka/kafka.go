package kafka

import (
	"fmt"
	"log"

	"github.com/Romasmi/advanced-url-shortener/internal/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConnection struct {
	Consumer *kafka.Consumer
	Producer *kafka.Producer
	Config   *config.Config
}

func (k *KafkaConnection) ConnectProducer() error {
	configMap := &kafka.ConfigMap{
		"bootstrap.severs": k.Config.Kafka.Brokers,
	}

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	k.Producer = producer
	go k.handleDeliveryReports()

	return nil
}

func (k *KafkaConnection) ConnectConsumer() error {
	configMap := &kafka.ConfigMap{
		"bootstrap.severs": k.Config.Kafka.Brokers,
	}
	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	k.Consumer = consumer
	return nil
}

func (k *KafkaConnection) Close() {
	if k.Producer != nil {
		k.Producer.Flush(15 * 1000)
		k.Producer.Close()
	}

	if k.Consumer != nil {
		err := k.Consumer.Close()
		if err != nil {
			log.Printf("error while closing Kafka consumer: %v", err)
		}
	}
}

func (k *KafkaConnection) handleDeliveryReports() {
	for e := range k.Producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
			} else {
				log.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}
