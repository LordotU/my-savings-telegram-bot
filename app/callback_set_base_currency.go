package app

import (
	"strings"

	"go.uber.org/zap"

	"github.com/LordotU/my-savings-telegram-bot/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *Application) botCallbackSetBaseCurrency(
	u *tgbotapi.Update,
	b *bot.Bot,
	params ...interface{},
) tgbotapi.Chattable {
	currency := params[0].(string)

	data := make(map[string]interface{})
	data["base"] = strings.ToUpper(currency)

	err := app.Repository.UpdateUser(u.CallbackQuery.From.ID, data)
	if err != nil {
		app.Logger.Error("Update user error", zap.Error(err))
		return nil
	}

	b.Api.AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "Base currency setted"))

	msg := tgbotapi.NewEditMessageText(
		u.CallbackQuery.Message.Chat.ID,
		u.CallbackQuery.Message.MessageID,
		"Now your base currency is *"+currency+"*",
	)
	msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{make([][]tgbotapi.InlineKeyboardButton, 0)}
	msg.ParseMode = tgbotapi.ModeMarkdown

	return msg
}
