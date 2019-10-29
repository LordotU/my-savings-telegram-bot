package app

import (
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/LordotU/my-savings-telegram-bot/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *Application) botHandlerAddSavings(
	u *tgbotapi.Update,
	b *bot.Bot,
	params ...interface{},
) tgbotapi.Chattable {
	commandArgs := u.Message.CommandArguments()
	isCommandArgsValid, err := regexp.Match("^\\d+([.,]\\d+)?\\s[a-zA-Z]{3}$", []byte(commandArgs))

	if commandArgs == "" || err != nil || !isCommandArgsValid {
		return tgbotapi.NewMessage(u.Message.Chat.ID, "Please, use command in following format:\n\n/add_savings {amount} {currency}")
	}

	commandArgsSplitted := strings.Split(commandArgs, " ")

	amount, err := strconv.ParseFloat(strings.Replace(commandArgsSplitted[0], ",", ".", 1), 64)
	if err != nil {
		app.Logger.Error("Create saving command args parsing error", zap.Any("args", commandArgsSplitted), zap.Error(err))
		return nil
	}

	currency := strings.ToUpper(commandArgsSplitted[1])
	if !strings.Contains(strings.Join(app.RatesProvider.GetSymbols(), ","), currency) {
		return tgbotapi.NewMessage(u.Message.Chat.ID, "Please, use command /get_currencies to get list of available currencies")
	}

	err = app.Repository.CreateSaving(u.Message.From.ID, amount, currency)
	if err != nil {
		app.Logger.Error("Create saving error", zap.Error(err))
		return nil
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "*"+strconv.FormatFloat(amount, 'f', -1, 64)+" "+currency+"* has been added!")
	msg.ParseMode = tgbotapi.ModeMarkdown

	return msg
}
