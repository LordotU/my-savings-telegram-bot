package providers

import (
	"encoding/json"

	"github.com/LordotU/my-savings-telegram-bot/repository"
	"github.com/LordotU/my-savings-telegram-bot/types"

	client "github.com/LordotU/go-fixerio"
	ratesTypes "github.com/LordotU/my-savings-telegram-bot/rates/types"
)

type FixerIOOptions struct {
	APIKey  string   `json:"FixerIOAPIKey"`
	Base    string   `json:"FixerIOBaseCurrency"`
	Secure  bool     `json:"FixerIOSecure"`
	Symbols []string `json:"FixerIOSymbols"`
}

type FixerIO struct {
	Client     *client.FixerIO
	Options    *FixerIOOptions
	Repository *repository.Repository
}

func NewFixerIO(o string, repository *repository.Repository) (*FixerIO, error) {
	options := &FixerIOOptions{}
	err := json.Unmarshal([]byte(o), options)
	if err != nil {
		return nil, err
	}

	client, err := client.New(options.APIKey, options.Base, options.Secure)
	if err != nil {
		return nil, err
	}

	return &FixerIO{
		Client:     client,
		Options:    options,
		Repository: repository,
	}, nil
}

func (f *FixerIO) UpdateRates() error {
	response, err := f.Client.GetLatest(f.Options.Symbols)
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

func (f *FixerIO) GetBaseCurrency() string {
	return f.Options.Base
}

func (f *FixerIO) GetSymbols() []string {
	return f.Options.Symbols
}
