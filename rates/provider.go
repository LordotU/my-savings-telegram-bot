package rates

import (
	"github.com/LordotU/my-savings-telegram-bot/repository"

	ratesProviders "github.com/LordotU/my-savings-telegram-bot/rates/providers"
	ratesTypes "github.com/LordotU/my-savings-telegram-bot/rates/types"
)

func GetNew(name string, settings ...interface{}) (provider ratesTypes.Provider, err error) {
	switch name {
	case "fixerio":
		provider, err = ratesProviders.GetNewFixerIO(
			settings[0].(string),
			settings[1].(string),
			settings[2].(bool),
			settings[3].(*repository.Repository),
		)
	default:
		provider, err = nil, nil
	}

	return
}
