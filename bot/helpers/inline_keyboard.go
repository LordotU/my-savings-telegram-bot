package helpers

import (
	"math"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type InlineKeyboardButtonFormatter func(item interface{}) string

func GetInlineKeyboard(
	btnsInRow int,
	items []interface{},
	buttonTextFormatter InlineKeyboardButtonFormatter,
	buttonDataFormatter InlineKeyboardButtonFormatter,
) *tgbotapi.InlineKeyboardMarkup {
	itemsLen := len(items)
	rowsLen := itemsLen / btnsInRow
	restBtnsLen := itemsLen % btnsInRow
	restRowsLen := 0
	if restBtnsLen > 0 {
		restRowsLen = int(math.Max(float64(btnsInRow-restBtnsLen), 0))
	}
	keyboardRows := make([][]tgbotapi.InlineKeyboardButton, rowsLen+restRowsLen)

	for i, item := range items {
		rowIdx := i / btnsInRow
		btnIdx := i % btnsInRow

		if btnIdx == 0 {
			keyboardRows[rowIdx] = make([]tgbotapi.InlineKeyboardButton, btnsInRow)

			if i == itemsLen-restBtnsLen {
				keyboardRows[rowIdx] = make([]tgbotapi.InlineKeyboardButton, restBtnsLen)
			}
		}

		keyboardRows[rowIdx][btnIdx] = tgbotapi.NewInlineKeyboardButtonData(
			buttonTextFormatter(item),
			buttonDataFormatter(item),
		)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	return &keyboard
}
