package app

import (
	"context"
	"errors"
	"math"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/LordotU/my-savings-telegram-bot/bot"
	"github.com/LordotU/my-savings-telegram-bot/rates"
	"github.com/LordotU/my-savings-telegram-bot/repository"

	ratesTypes "github.com/LordotU/my-savings-telegram-bot/rates/types"
)

type Application struct {
	Bot           *bot.Bot
	Config        *Config
	Ctx           context.Context
	CtxCancel     context.CancelFunc
	Error         chan error
	Logger        *zap.Logger
	MongoDB       *mongo.Client
	RatesProvider ratesTypes.Provider
	Repository    *repository.Repository
	WaitGroup     *sync.WaitGroup
}

func New() (*Application, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	logger, err := GetLogger(config.Debug)
	if err != nil {
		return nil, err
	}

	app := &Application{
		Config:    config,
		Error:     make(chan error, math.MaxUint8),
		Logger:    logger,
		WaitGroup: new(sync.WaitGroup),
	}
	app.Ctx, app.CtxCancel = context.WithCancel(context.Background())

	defer func() {
		if err != nil {
			app.Close()
		}
	}()

	app.initMongoDB()
	app.initRepository()
	app.initBot()
	app.initRatesProvider()

	logger.Info("Application created")

	return app, nil
}

func (app *Application) Run() {
	var err error
	defer app.Close()

	app.Bot.AddCommandHandler("default", app.botHandlerDefault)
	app.Bot.AddCommandHandler("help", app.botHandlerHelp)
	app.Bot.AddCommandHandler("get_currencies", app.botHandlerGetCurrencies)
	app.Bot.AddCommandHandler("register", app.middlewareGetExistingUser(app.botHandlerRegister, false, true))
	app.Bot.AddCommandHandler("set_base_currency", app.middlewareGetExistingUser(app.botHandlerSetBaseCurrency, true, false))
	app.Bot.AddCallbackHandler("set_base_currency", app.botCallbackSetBaseCurrency)
	app.Bot.AddCommandHandler("add_savings", app.middlewareGetExistingUser(app.botHandlerAddSavings, true, false))
	app.Bot.AddCommandHandler("get_savings", app.middlewareGetExistingUser(app.botHandlerGetSavings, true, false))
	app.Bot.AddCallbackHandler("remove_savings", app.botCallbackRemoveSavings)

	go app.runTelegramBot()
	app.runRatesUpdates()

	app.Logger.Info("Application started")

	select {
	case err = <-app.Error:
		app.Logger.Panic("Application error", zap.Error(err))
	case <-app.Ctx.Done():
		app.Logger.Error("Application stops via Ctx")
	case signal := <-waitForExitSignal():
		app.Logger.Warn("Application was interrupted", zap.Stringer("signal", signal))
	}
}

func (app *Application) Stop() {
	app.Logger.Warn("Application trying to stop")
	app.CtxCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	go func() {
		defer cancel()
		app.WaitGroup.Wait()
	}()

	<-ctx.Done()

	if ctx.Err() != context.Canceled {
		app.Logger.Panic("Application stopped by timeout")
	} else {
		app.Logger.Warn("Application successfully stoppped")
	}
}

func (app *Application) Close() {
	defer app.Stop()
	app.Logger.Warn("Application trying to close open connections")

	defer func() {
		if app.MongoDB != nil {
			if err := app.MongoDB.Disconnect(app.Ctx); err != nil {
				app.Logger.Warn("MongoDB disconnecting error", zap.Error(err))
			} else {
				app.Logger.Info("MongoDB disconnecting success")
			}
		}
	}()

	defer func() {
		if app.Bot != nil {
			app.Bot.Api.StopReceivingUpdates()
			app.Logger.Info("TelegramBotAPI disconnecting success")
		}
	}()
}

func (app *Application) initMongoDB() {
	var err error

	app.MongoDB, err = mongo.Connect(app.Ctx, options.Client().ApplyURI(app.Config.MongoDbURI))
	if err != nil {
		app.Logger.Panic(
			"Cannot init MongoDB connection",
			zap.Error(err),
			zap.String("MongoDbURI", app.Config.MongoDbURI),
		)
	}

	err = app.MongoDB.Ping(app.Ctx, nil)
	if err != nil {
		app.Logger.Panic(
			"Cannot ping MongoDB connection",
			zap.Error(err),
			zap.String("MongoDbURI", app.Config.MongoDbURI),
		)
	}
}

func (app *Application) initRepository() {
	app.Repository = repository.GetNew(app.MongoDB.Database(app.Config.MongoDbDatabase))
}

func (app *Application) initBot() {
	var err error
	app.Bot, err = bot.GetNew(app.Config.TelegramAPIToken, app.Config.DebugTelegramAPI, app.Logger)
	if err != nil {
		app.Logger.Error("Cannot init TelegramBotAPI", zap.Error(err))
	}
}

func (app *Application) initRatesProvider() {
	var err error
	if app.RatesProvider, err = rates.GetNew(
		app.Config.RatesProvider,
		app.Config.FixerIOAPIKey,
		app.Config.FixerIOBaseCurrency,
		app.Config.FixerIOSecure,
		app.Repository,
	); err != nil {
		app.Logger.Panic("Cannot init RatesProvider", zap.Error(err))
	}
}

func (app *Application) runTelegramBot() {
	err := app.Bot.Run(0, app.Config.TelegramAPIUpdatesTimeout)
	if err != nil {
		app.Logger.Panic("Cannot run Telegram Bot", zap.Error(err))
	}
}

func (app *Application) getTimerForRatesUpdates(duration int) *time.Timer {
	return time.NewTimer(time.Duration(duration) * time.Second)
}

func (app *Application) runRatesUpdates() {
	timer := app.getTimerForRatesUpdates(0)

	go func() {
		app.WaitGroup.Add(1)
		defer func() {
			app.WaitGroup.Done()
		}()

		app.Logger.Info("Start rates update")

		for {
			select {
			case <-app.Ctx.Done():
				return
			case _, ok := <-timer.C:
				if ok {
					app.Logger.Info("Process rates update")

					if err := app.RatesProvider.UpdateRates(strings.Split(app.Config.FixerIOSymbols, ",")); err != nil {
						app.Logger.Error("Process rates error", zap.Error(err))
					}

					timer = app.getTimerForRatesUpdates(app.Config.ExchangeRatesGettingTimer)
				} else {
					app.Error <- errors.New("runRatesUpdates timer stops")
					return
				}
			}
		}
	}()
}

func waitForExitSignal() (signalsCh chan os.Signal) {
	signalsCh = make(chan os.Signal, 1)
	signal.Notify(signalsCh, os.Interrupt, syscall.SIGTERM)

	return signalsCh
}
