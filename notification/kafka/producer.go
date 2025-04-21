package kafka

import (
	"encoding/json"
	"p2platform/notification/model"

	"github.com/IBM/sarama"
)

func Publish(producer sarama.SyncProducer, topic string, msg model.NotifictationMessage) error {
	value, _ := json.Marshal(msg)
	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	_, _, err := producer.SendMessage(kafkaMsg)
	return err
}
