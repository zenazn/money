package money

import (
	"github.com/zenazn/money/currency"
	"github.com/zenazn/money/decimal"
)

type ExchangeRate struct {
	Source      currency.Currency
	Destination currency.Currency
	Rate        decimal.Rate
}

func (e ExchangeRate) check() {
	if e.Source == nil {
		panic("money: exchange rate source currency must be non-nil")
	}
	if e.Destination == nil {
		panic("money: exchange rate destination currency must be non-nil")
	}
}

// ExchangeE performs a currency exchange calculation, returning the converted
// amount in the destination currency, or an error if the source amount does not
// match the source currency on the exchange rate.
func (m Money) ExchangeE(e ExchangeRate) (Money, error) {
	e.check()
	if err := compat(m.ccy, e.Source); err != nil {
		return Money{}, err
	}
	return Money{m.amt.Mul(e.Rate), e.Destination}, nil
}

// Exchange performs a currency exchange calculation, returning the converted
// amount in the destination currency. If the source currency
func (m Money) Exchange(e ExchangeRate) Money {
	e.check()
	if err := compat(m.ccy, e.Source); err != nil {
		panic(err.Error())
	}

	return Money{m.amt.Mul(e.Rate), e.Destination}
}
