package types

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type User struct {
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`

	ID           int    `bson:"telegram_id"`
	FirstName    string `bson:"first_name"`
	LastName     string `bson:"last_name"`
	UserName     string `bson:"username"`
	LanguageCode string `bson:"language_code"`

	BaseCurrency string `bson:"base"`
}

func GetNewUser(tgUser tgbotapi.User, baseCurrency string) *User {
	return &User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		ID:           tgUser.ID,
		FirstName:    tgUser.FirstName,
		LastName:     tgUser.LastName,
		UserName:     tgUser.UserName,
		LanguageCode: tgUser.LanguageCode,

		BaseCurrency: baseCurrency,
	}
}
