package token

type Payload struct {
	TelegramId int64  `json:"telegram_id"`
	Username   string `json:"username"`
}

func NewPayload(telegramId int64, username string) (*Payload, error) {
	payload := &Payload{
		TelegramId: telegramId,
		Username:   username,
	}
	return payload, nil
}
