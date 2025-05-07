package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBSource          string `mapstructure:"DB_SOURCE"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	Environment       string `mapstructure:"ENVIRONMENT"`
	TelegramBotToken  string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	KafkaBrokers      string `mapstructure:"KAFKA_BROKERS"`
	AccessTokenDuration string `mapstructure:"ACCESS_TOKEN_DURATION"`
	BaseURL           string `mapstructure:"BASE_URL"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
