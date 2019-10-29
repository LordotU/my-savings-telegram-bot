package app

import "github.com/caarlos0/env/v6"

type Config struct {
	Debug                     bool   `env:"DEBUG" envDefault:"false"`
	DebugTelegramAPI          bool   `env:"DEBUG_TELEGRAM_API" envDefault:"false"`
	MongoDbURI                string `env:"MONGO_DB_URI" envDefault:"mongodb://localhost:27017"`
	MongoDbDatabase           string `env:"MONGO_DB_DATABASE" envDefault:"my-savings-telegram-bot"`
	RatesProvider             string `env:"RATES_PROVIDER" envDefault:"fixerio"`
	RatesProviderOptions      string `env:"RATES_PROVIDER_OPTIONS" envDefault:"{\"FixerIOAPIKey\": \"\", \"FixerIOBaseCurrency\": \"EUR\", \"FixerIOSecure\": false, \"FixerIOSymbols\": [\"USD\",\"EUR\",\"JPY\",\"GBP\",\"AUD\",\"CAD\",\"CHF\",\"CNY\",\"NZD\",\"RUB\"]}"`
	RatesProviderUpdatePeriod int    `env:"RATES_RPOVIDER_UPDATE_PERIOD" envDefault:"3600"`
	TelegramAPIToken          string `env:"TELEGRAM_API_TOKEN" envDefault:""`
	TelegramAPIUpdatesTimeout int    `env:"TELEGRAM_API_UPDATES_TIMEOUT" envDfault:"60"`
}

func GetConfig() (config *Config, err error) {
	config = new(Config)
	return config, env.Parse(config)
}
