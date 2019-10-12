package bot

import (
	"errors"
	"strings"

	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Handler func(update *tgbotapi.Update, bot *Bot, params ...interface{}) tgbotapi.Chattable

type Bot struct {
	Api              *tgbotapi.BotAPI
	CallbacksHandler map[string]Handler
	CommandsHandlers map[string]Handler
	Logger           *zap.Logger
}

func GetNew(apiToken string, debug bool, logger *zap.Logger) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, err
	}

	api.Debug = debug

	return &Bot{
		Api:              api,
		CallbacksHandler: make(map[string]Handler),
		CommandsHandlers: make(map[string]Handler),
		Logger:           logger,
	}, nil
}

func (bot *Bot) AddCommandHandler(command string, handler Handler) {
	bot.CommandsHandlers[command] = handler
}

func (bot *Bot) AddCallbackHandler(command string, handler Handler) {
	bot.CallbacksHandler[command] = handler
}

func (bot *Bot) Run(offset, timeout int) error {
	u := tgbotapi.NewUpdate(offset)
	u.Timeout = timeout

	updates, err := bot.Api.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	updates.Clear()

	for update := range updates {
		switch true {
		case update.Message != nil && update.Message.IsCommand():
			go bot.ProcessCommand(&update)
		case update.CallbackQuery != nil:
			go bot.ProcessCallbackQuery(&update)
		default:
			continue
		}
	}

	return nil
}

func (bot *Bot) ProcessCommand(u *tgbotapi.Update) {
	bot.Logger.Info(
		"Incoming command message",
		zap.String("from", u.Message.From.UserName),
		zap.String("text", u.Message.Text),
	)

	command := u.Message.Command()

	msg := bot.CommandsHandlers["default"](u, bot)
	if handler, ok := bot.CommandsHandlers[command]; ok {
		msg = handler(u, bot)
	}

	bot.sendMessage(msg)
}

func (bot *Bot) ProcessCallbackQuery(u *tgbotapi.Update) {
	bot.Logger.Info(
		"Incoming callback query",
		zap.String("from", u.CallbackQuery.Message.From.UserName),
		zap.String("data", u.CallbackQuery.Data),
	)

	query := strings.Split(u.CallbackQuery.Data, "/")

	var msg tgbotapi.Chattable

	if handler, ok := bot.CallbacksHandler[query[0]]; ok {
		msg = handler(u, bot, query[1])
	}

	bot.sendMessage(msg)
}

func (bot *Bot) sendMessage(msg tgbotapi.Chattable) {
	if msg == nil {
		bot.Logger.Error("Cannot send msg", zap.Error(errors.New("msg cannot be nil")))
		return
	}

	if _, err := bot.Api.Send(msg); err != nil {
		bot.Logger.Error("Cannot send msg", zap.Error(err))
		return
	}
}
