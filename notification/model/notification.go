package model

type NotifictationMessage struct {
	TelegramId int64  `json:"telegram_id"`
	Message    string `json:"message"`
	EventType  string `json:"event_type"`
}
