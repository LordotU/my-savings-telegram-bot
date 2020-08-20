package app

import (
	"strconv"

	"go.uber.org/zap"

	"github.com/LordotU/my-savings-telegram-bot/bot"
	"github.com/LordotU/my-savings-telegram-bot/types"

	botHelpers "github.com/LordotU/my-savings-telegram-bot/bot/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *Application) botHandlerGetSavings(
	u *tgbotapi.Update,
	b *bot.Bot,
	params ...interface{},
) tgbotapi.Chattable {
	user := params[0].(*types.User)

	savings, err := app.Repository.FindSavings(u.Message.From.ID)
	if err != nil && err.Error() != "mongo: no documents in result" {
		app.Logger.Error("Find savings error", zap.Error(err))
		return nil
	}

	if len(savings) == 0 {
		return tgbotapi.NewMessage(u.Message.Chat.ID, "You have no savings now. Please, use command:\n\n/add_savings {amount} {currency}")
	}

	userBaseCurrencyRate, err := app.Repository.FindRate(user.BaseCurrency)
	if err != nil {
		app.Logger.Error("Find rate error", zap.Error(err))
		return nil
	}

	savingsInterface := make([]interface{}, len(savings))
	ratesInterface := make([]interface{}, len(savings))
	totalInUserBaseCurrency := float64(0)
	for i, saving := range savings {
		savingsInterface[i] = (*saving)
		rate, err := app.Repository.FindRate(saving.Currency)
		if err != nil {
			app.Logger.Error("Get total savings in user base currency error", zap.Error(err))
			return nil
		}
		ratesInterface[i] = struct {
			Currency string
			Rate     string
		}{
			Currency: saving.Currency,
			Rate:     strconv.FormatFloat(userBaseCurrencyRate.Rate / rate.Rate, 'f', -1, 32),
		}
		totalInUserBaseCurrency += saving.Amount / rate.Rate * userBaseCurrencyRate.Rate
	}

	msg, err := botHelpers.GetMsgFromMdTemplate(
		"get_savings.md",
		struct {
			TotalInUserBaseCurrency string
			UserBaseCurrency        string
			Rates                   []interface{}
		}{
			strconv.FormatFloat(totalInUserBaseCurrency, 'f', -1, 32),
			userBaseCurrencyRate.Currency,
			ratesInterface,
		},
		u.Message.Chat.ID,
	)
	if err != nil {
		app.Logger.Error(
			"Cannot process template file",
			zap.String("file", "templates/get_savings.md"),
			zap.Error(err),
		)
		return nil
	}

	msgConfig := msg.(tgbotapi.MessageConfig)
	msgConfig.ReplyMarkup = botHelpers.GetInlineKeyboard(
		1,
		savingsInterface,
		func(item interface{}) string {
			return strconv.FormatFloat(item.(types.Saving).Amount, 'f', -1, 32) + " " + item.(types.Saving).Currency
		},
		func(item interface{}) string {
			return "remove_savings/" + item.(types.Saving).ID.Hex()
		},
	)

	return msgConfig
}
