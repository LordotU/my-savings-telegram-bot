package app

import "github.com/caarlos0/env/v6"

type Config struct {
	Debug                     bool   `env:"DEBUG" envDefault:"false"`
	DebugTelegramAPI          bool   `env:"DEBUG_TELEGRAM_API" envDefault:"false"`
	ExchangeRatesGettingTimer int    `env:"EXCHANGE_RATES_GETTING_TIMER" envDefault:"3600"`
	FixerIOAPIKey             string `env:"FIXERIO_API_KEY" envDefault:""`
	FixerIOBaseCurrency       string `env:"FIXERIO_BASE_CURRENCY" envDefault:"EUR"`
	FixerIOSecure             bool   `env:"FIXERIO_SECURE" envDefault:"false"`
	FixerIOSymbols            string `env:"FIXERIO_SYMBOLS" envDefault:"USD,EUR,JPY,GBP,AUD,CAD,CHF,CNY,NZD,RUB"`
	MongoDbURI                string `env:"MONGO_DB_URI" envDefault:"mongodb://localhost:27017"`
	MongoDbDatabase           string `env:"MONGO_DB_DATABASE" envDefault:"my-savings-telegram-bot"`
	RatesProvider             string `env:"RATES_PROVIDER" envDefault:"fixerio"`
	TelegramAPIToken          string `env:"TELEGRAM_API_TOKEN" envDefault:""`
	TelegramAPIUpdatesTimeout int    `env:"TELEGRAM_API_UPDATES_TIMEOUT" envDfault:"60"`
}

func GetConfig() (config *Config, err error) {
	config = new(Config)
	return config, env.Parse(config)
}
