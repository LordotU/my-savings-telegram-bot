package app

import (
	"go.uber.org/zap"

	"github.com/LordotU/my-savings-telegram-bot/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *Application) botCallbackRemoveSavings(
	u *tgbotapi.Update,
	b *bot.Bot,
	params ...interface{},
) tgbotapi.Chattable {
	savingsIDs := make([]string, 1)
	if params[0] != nil && params[0].(string) != "" {
		savingsIDs[0] = params[0].(string)
	}

	err := app.Repository.DeleteSavings(savingsIDs)
	if err != nil {
		app.Logger.Error("Delete saving error", zap.Error(err))
		return nil
	}

	b.Api.AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "Saving deleted"))

	msg := tgbotapi.NewEditMessageText(
		u.CallbackQuery.Message.Chat.ID,
		u.CallbackQuery.Message.MessageID,
		"Please, use /get_savings command again to get updated list",
	)
	msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{make([][]tgbotapi.InlineKeyboardButton, 0)}

	return msg
}
