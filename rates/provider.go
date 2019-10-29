package rates

import (
	"github.com/LordotU/my-savings-telegram-bot/repository"

	ratesProviders "github.com/LordotU/my-savings-telegram-bot/rates/providers"
	ratesTypes "github.com/LordotU/my-savings-telegram-bot/rates/types"
)

func New(name string, options string, repository *repository.Repository) (provider ratesTypes.Provider, err error) {
	switch name {
	case "fixerio":
		provider, err = ratesProviders.NewFixerIO(
			options,
			repository,
		)
	default:
		provider, err = nil, nil
	}

	return
}
