package types

type CurrencyExchangeRate struct {
	Currency string
	Base     string
	Rate     float64
}

type Provider interface {
	GetBaseCurrency() string
	GetSymbols() []string
	UpdateRates() error
}
