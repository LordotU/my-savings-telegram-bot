package app

import (
	"go.uber.org/zap"

	"github.com/LordotU/my-savings-telegram-bot/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *Application) middlewareGetExistingUser(
	h bot.Handler,
	preventNegative bool,
	preventPositive bool,
) bot.Handler {
	return func(
		u *tgbotapi.Update,
		b *bot.Bot,
		params ...interface{},
	) tgbotapi.Chattable {
		existingUser, err := app.Repository.FindUser(u.Message.From.ID)
		if err != nil && err.Error() != "mongo: no documents in result" {
			app.Logger.Error("Find of existing user error", zap.Error(err))
			return nil
		}

		if preventNegative && (existingUser == nil || existingUser.ID == 0) {
			return tgbotapi.NewMessage(u.Message.Chat.ID, "Please, use /register command firstly")
		}

		if preventPositive && existingUser != nil && existingUser.ID != 0 {
			return tgbotapi.NewMessage(u.Message.Chat.ID, "You are already registered!")
		}

		return h(u, b, existingUser)
	}
}
