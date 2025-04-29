package token

type Maker interface {
	CreateToken(telegramId int64, username string) (string, error)
	VerifyToken(token string) (*Payload, error)
}
