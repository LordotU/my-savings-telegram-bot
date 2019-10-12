package app

import (
	"go.uber.org/zap"

	"github.com/LordotU/my-savings-telegram-bot/bot"
	"github.com/LordotU/my-savings-telegram-bot/types"

	botHelpers "github.com/LordotU/my-savings-telegram-bot/bot/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *Application) botHandlerSetBaseCurrency(
	u *tgbotapi.Update,
	b *bot.Bot,
	params ...interface{},
) tgbotapi.Chattable {
	rates, err := app.Repository.FindRates([]string{})
	if err != nil && err.Error() != "mongo: no documents in result" {
		app.Logger.Error("Find rates error", zap.Error(err))
		return nil
	}

	ratesInterface := make([]interface{}, len(rates))
	for i, rate := range rates {
		ratesInterface[i] = (*rate)
	}

	keyboard := botHelpers.GetInlineKeyboard(
		3,
		ratesInterface,
		func(item interface{}) string {
			return item.(types.Rate).Currency
		},
		func(item interface{}) string {
			return "set_base_currency/" + item.(types.Rate).Currency
		},
	)

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Set one of suggested currencies as your base:")
	msg.ReplyMarkup = keyboard

	return msg
}
