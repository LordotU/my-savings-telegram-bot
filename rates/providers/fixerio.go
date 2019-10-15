package providers

import (
	"github.com/LordotU/my-savings-telegram-bot/repository"
	"github.com/LordotU/my-savings-telegram-bot/types"

	client "github.com/LordotU/go-fixerio"
	ratesTypes "github.com/LordotU/my-savings-telegram-bot/rates/types"
)

type FixerIO struct {
	Client     *client.FixerIO
	Repository *repository.Repository
}

func GetNewFixerIO(APIKey string, Base string, Secure bool, repository *repository.Repository) (*FixerIO, error) {
	f, err := client.New(APIKey, Base, Secure)
	if err != nil {
		return nil, err
	}

	return &FixerIO{
		f,
		repository,
	}, nil
}

func (f *FixerIO) UpdateRates(params ...interface{}) error {
	response, err := f.Client.GetLatest(params[0].([]string))
	if err != nil {
		return err
	}

	result := make([]*types.Rate, 0, len(response.Rates))
	for key, value := range response.Rates {
		result = append(result, types.GetNewRate(ratesTypes.CurrencyExchangeRate{
			Currency: key,
			Base:     f.Client.Base,
			Rate:     value,
		}))
	}

	return f.Repository.UpdateRates(result)
}
