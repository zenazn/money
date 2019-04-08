package money

import (
	"github.com/zenazn/money/currency"
	"github.com/zenazn/money/decimal"
)

// ExchangeRate describes a Rate at which the Source currency can be sold to buy
// the Destination currency.
type ExchangeRate struct {
	src, dst currency.Currency
	rate     decimal.Rate
}

// Source returns the source currency of the exchange rate.
func (e ExchangeRate) Source() currency.Currency {
	return e.src
}

// Destination returns the destination currency of the exchange rate.
func (e ExchangeRate) Destination() currency.Currency {
	return e.dst
}

// Rate returns the exchange rate as a scalar.
func (e ExchangeRate) Rate() decimal.Rate {
	return e.rate
}

// NewExchangeRate returns an exchange rate with the given information. It
// panics if either of the currencies is nil.
func NewExchangeRate(src, dst currency.Currency, rate decimal.Rate) ExchangeRate {
	if src == nil {
		panic("money: exchange rate source currency must be non-nil")
	}
	if dst == nil {
		panic("money: exchange rate destination currency must be non-nil")
	}
	return ExchangeRate{src, dst, rate}
}

// ExchangeErr performs a currency exchange calculation, returning the converted
// amount in the destination currency, or an error if the source amount does not
// match the source currency on the exchange rate.
func (m Money) ExchangeErr(e ExchangeRate) (Money, error) {
	if err := compat(m.ccy, e.src); err != nil {
		return Money{}, err
	}

	// TODO: would be neat to be able to fuse these ops to prevent loss of
	// precision
	d := m.amt.Mul(e.rate)

	// If the scaling factors differ, we'll need to scale the value to
	// adjust. This ideally shouldn't happen much, so for simplicity the
	// code below is naive / slower than it needs to be.
	su := e.src.Units()
	du := e.dst.Units()
	if su.MajorUnitScalingFactorExponent > du.MajorUnitScalingFactorExponent {
		sfd := int(su.MajorUnitScalingFactorExponent - du.MajorUnitScalingFactorExponent)
		for i := 0; i < sfd; i++ {
			d = d.Div(decimal.NewRate(10 * 1000000))
		}

	} else if su.MajorUnitScalingFactorExponent < du.MajorUnitScalingFactorExponent {
		sfd := int(du.MajorUnitScalingFactorExponent - su.MajorUnitScalingFactorExponent)
		for i := 0; i < sfd; i++ {
			d = d.Mul(decimal.NewRate(10 * 1000000))
		}
	}

	return Money{d, e.dst}, nil
}

// Exchange performs a currency exchange calculation, returning the converted
// amount in the destination currency. This function panics if the source amount
// does not match the source currency on the exchange rate.
func (m Money) Exchange(e ExchangeRate) Money {
	m2, err := m.ExchangeErr(e)
	if err != nil {
		panic(err.Error())
	}
	return m2
}
