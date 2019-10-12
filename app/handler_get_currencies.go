package app

import (
	"strings"

	"go.uber.org/zap"

	"github.com/LordotU/my-savings-telegram-bot/bot"

	botHelpers "github.com/LordotU/my-savings-telegram-bot/bot/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *Application) botHandlerGetCurrencies(
	u *tgbotapi.Update,
	b *bot.Bot,
	params ...interface{},
) tgbotapi.Chattable {
	rates, err := app.Repository.FindRates([]string{})
	if err != nil && err.Error() != "mongo: no documents in result" {
		app.Logger.Error("Find rates error", zap.Error(err))
		return nil
	}

	currencies := make([]string, len(rates))
	for i, rate := range rates {
		currencies[i] = rate.Currency
	}

	msg, err := botHelpers.GetMsgFromMdTemplate(
		"get_currencies.md",
		struct{ Currencies string }{strings.Join(currencies, ", ")},
		u.Message.Chat.ID,
	)
	if err != nil {
		app.Logger.Error(
			"Cannot process template file",
			zap.String("file", "templates/get_currencies.md"),
			zap.Error(err),
		)
		return nil
	}

	return msg
}
