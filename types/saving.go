package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Saving struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`

	TelegramID int     `bson:"telegram_id"`
	Amount     float64 `bson:"amount"`
	Currency   string  `bson:"currency"`
}

func GetNewSaving(telegramID int, amount float64, currency string) *Saving {
	return &Saving{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		TelegramID: telegramID,
		Amount:     amount,
		Currency:   currency,
	}
}
