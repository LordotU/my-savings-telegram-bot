package app

import (
	"go.uber.org/zap"

	"github.com/LordotU/my-savings-telegram-bot/bot"

	botHelpers "github.com/LordotU/my-savings-telegram-bot/bot/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *Application) botHandlerRegister(
	u *tgbotapi.Update,
	b *bot.Bot,
	params ...interface{},
) tgbotapi.Chattable {
	_, err := app.Repository.CreateUser(*u.Message.From, app.RatesProvider.GetBaseCurrency())
	if err != nil {
		app.Logger.Error("Create user error", zap.Error(err))
		return nil
	}

	msg, err := botHelpers.GetMsgFromMdTemplate(
		"register.md",
		struct{ BaseCurrency string }{app.RatesProvider.GetBaseCurrency()},
		u.Message.Chat.ID,
	)
	if err != nil {
		app.Logger.Error(
			"Cannot process template file",
			zap.String("file", "templates/register.md"),
			zap.Error(err),
		)
		return nil
	}

	return msg
}
