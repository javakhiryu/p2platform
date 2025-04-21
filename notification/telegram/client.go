package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	Bot *tgbotapi.BotAPI
}

func New(botToken string) *Client{
	bot, err :=tgbotapi.NewBotAPI(botToken)
	if err !=nil {
		panic(err)
	}
	return &Client{
		Bot: bot,
	}
}
func (c *Client) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := c.Bot.Send(msg)
	return err
}