package helpers

import (
	"github.com/LordotU/my-savings-telegram-bot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetMsgFromMdTemplate(
	templateName string,
	vars interface{},
	chatID int64,
) (tgbotapi.Chattable, error) {
	processedTemplate, err := utils.ProcessTemplateFile(
		"templates/"+templateName,
		vars,
	)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(chatID, processedTemplate)
	msg.ParseMode = tgbotapi.ModeMarkdown

	return msg, nil
}
