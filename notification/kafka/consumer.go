package kafka

import (
	"encoding/json"
	"fmt"
	"p2platform/notification/model"
	"p2platform/notification/telegram"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

func StartConsumer(brokers []string, topic string, tg *telegram.Client) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err !=nil {
		return err
	}
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err !=nil{
		return err
	} 
	log.Info().Msg("Kafka consumer started...")
	for msg :=range partitionConsumer.Messages() {
		var notif model.NotifictationMessage
		if err := json.Unmarshal(msg.Value, &notif); err !=nil {
			log.Error().Str("Invalid message format:", err.Error())
			continue
		}
		//log.Info().Msg(fmt.Sprintf("Received from %d: %s", notif.TelegramId, notif.Message))
		if err :=tg.SendMessage(notif.TelegramId, notif.Message); err !=nil{
			log.Error().Str("Telegram send error:", err.Error())
		} else {
			log.Info().Msg(fmt.Sprintf("Sent to %d: %s", notif.TelegramId, notif.Message))
		}
	}
	return nil
}