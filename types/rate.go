package types

import (
	"time"

	ratesTypes "github.com/LordotU/my-savings-telegram-bot/rates/types"
)

type Rate struct {
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`

	Currency string  `bson:"_id"`
	Base     string  `bson:"base"`
	Rate     float64 `bson:"rate"`
}

func GetNewRate(rate ratesTypes.CurrencyExchangeRate) *Rate {
	return &Rate{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		Currency: rate.Currency,
		Base:     rate.Base,
		Rate:     rate.Rate,
	}
}
