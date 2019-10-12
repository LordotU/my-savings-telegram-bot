package types

type CurrencyExchangeRate struct {
	Currency string
	Base     string
	Rate     float64
}

type Provider interface {
	UpdateRates(params ...interface{}) error
}
