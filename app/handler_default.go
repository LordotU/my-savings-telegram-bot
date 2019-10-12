package app

import (
	"go.uber.org/zap"

	"github.com/LordotU/my-savings-telegram-bot/bot"

	botHelpers "github.com/LordotU/my-savings-telegram-bot/bot/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *Application) botHandlerDefault(
	u *tgbotapi.Update,
	b *bot.Bot,
	params ...interface{},
) tgbotapi.Chattable {
	msg, err := botHelpers.GetMsgFromMdTemplate(
		"default.md",
		struct{}{},
		u.Message.Chat.ID,
	)
	if err != nil {
		app.Logger.Error(
			"Cannot process template file",
			zap.String("file", "templates/default.md"),
			zap.Error(err),
		)
		return nil
	}

	return msg
}
