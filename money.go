package money

import (
	"fmt"
	"strings"

	"github.com/zenazn/money/currency"
	"github.com/zenazn/money/decimal"
)

// Money represents an amount of a currency. It stores the amount as a fixed
// point decimal value, so it is appropriate for use as a data type for
// accounting. Money values of different currencies are incompatible.
//
// The zero value of Money is a "currencyless zero" and may be immediately used.
// Unlike other values, a currencyless zero is compatible with values of all
// currencies, and behaves as if it represented zero units of the currency of
// the other operand.
type Money struct {
	amt decimal.Decimal
	ccy currency.Currency
}

// Zero returns the zero value of the given non-nil currency.
func Zero(ccy currency.Currency) Money {
	if ccy == nil {
		panic("money: no currency given")
	}
	return Money{ccy: ccy}
}

// New returns a Money with the given amount and non-nil currency.
func New(amt decimal.Decimal, ccy currency.Currency) Money {
	if ccy == nil {
		panic("money: no currency given")
	}
	return Money{amt, ccy}
}

// FromMinorCurrencyUnits returns a Money with the given currency and amount, as
// represented as an integer number of minor currency units in that currency.
func FromMinorCurrencyUnits(amt int64, ccy currency.Currency) Money {
	if ccy == nil {
		panic("money: no currency given")
	}
	sf := 1000000
	u := ccy.Units()
	sfe := int(u.MajorUnitScalingFactorExponent - u.MinorUnitsInMajorUnitExponent)
	for i := 0; i < sfe; i++ {
		sf = sf * 10
	}
	d := decimal.FromI64(amt).Mul(decimal.NewRate(sf))
	return Money{d, ccy}
}

// Parse interprets the amount as a floating-point decimal string (using "." to
// separate the whole part from the fractional part, and without thousands
// separators or other adornments) and the currency as an ISO 4217 currency
// code, and returns a Money representing that value, or an error if either the
// amount or currency fails to parse.
func Parse(amt, ccy string) (Money, error) {
	// TODO
	return Money{}, nil
}

// Amount returns a decimal integer number of minimum-representable-units of the
// value. In order to interpret the relationship between a currency's
// minimum-representable-unit and major units, please consult the currency's
// Units.
func (m Money) Amount() decimal.Decimal {
	return m.amt
}

// Currency returns the currency of this value. As a special case, the special
// "currencyless zero" value will return a nil currency.
func (m Money) Currency() currency.Currency {
	return m.ccy
}

// Zero returns true if this value is zero.
func (m Money) Zero() bool {
	return m.amt == decimal.Decimal{}
}

func compat(a, b currency.Currency) error {
	// Fast path: currencies are identical objects
	if a == b {
		return nil
	}
	// Special case: currencyless zeroes are compatible
	if a == nil || b == nil {
		return nil
	}
	// Slower path: currencies are identical if they have identical symbols
	as := a.Symbol()
	bs := b.Symbol()
	if as == bs {
		return nil
	}
	return fmt.Errorf("money: incompatible currencies %s and %s", as, bs)
}

func (m Money) compatCcy(o Money) currency.Currency {
	// Operations on currencyless zeroes are "sticky"
	if m.ccy != nil {
		return m.ccy
	}
	return o.ccy

}

// AddE adds the two values and returns the result, or an error if the two
// values were incompatible.
func (m Money) AddE(o Money) (Money, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return Money{}, err
	}
	return Money{m.amt.Add(o.amt), m.compatCcy(o)}, nil
}

// Add adds the two values and returns the result. If the two values are
// incompatible, Add will panic.
func (m Money) Add(o Money) Money {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return Money{m.amt.Add(o.amt), m.compatCcy(o)}
}

// SubE subtracts the second value from the first and returns the result, or an
// error if the two values are incompatible.
func (m Money) SubE(o Money) (Money, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return Money{}, err
	}
	return Money{m.amt.Sub(o.amt), m.compatCcy(o)}, nil
}

// Sub subtracts the second value from the first value and returns the result.
// If the two values are incompatible, Sub will panic.
func (m Money) Sub(o Money) Money {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return Money{m.amt.Sub(o.amt), m.compatCcy(o)}
}

// Neg negates the value and returns the result.
func (m Money) Neg() Money {
	return Money{m.amt.Neg(), m.ccy}
}

// Mul multiplies the value by the given scalar rate and returns the result.
func (m Money) Mul(r decimal.Rate) Money {
	return Money{m.amt.Mul(r), m.ccy}
}

// Div divides the value by the given scalar rate and returns the result.
func (m Money) Div(r decimal.Rate) Money {
	return Money{m.amt.Div(r), m.ccy}
}

func (m Money) RoundToMinorUnits() Money {
	// TODO
	return m
}

// String renders the amount as a human-readable currency and amount without
// loss of precision, like "EUR 1.30", "JPY 990", or "USD 0.0187".
func (m Money) String() string {
	if m.ccy == nil {
		return "0"
	}

	s := m.amt.String()
	var prefix string
	if s[0] == '-' {
		s = s[1:]
		prefix = m.ccy.Symbol() + " -"
	} else {
		prefix = m.ccy.Symbol() + " "
	}

	u := m.ccy.Units()
	sf := int(u.MajorUnitScalingFactorExponent)

	if len(s) > sf {
		s = s[:len(s)-sf] + "." + s[len(s)-sf:]
	} else if len(s) == sf {
		s = "0." + s
	} else {
		s = "0." + strings.Repeat("0", sf-len(s)) + s
	}

	walkback := sf - int(u.MinorUnitsInMajorUnitExponent)
	// Special case: for currencies without minor units, also consume the
	// decimal point
	if u.MinorUnitsInMajorUnitExponent == 0 {
		walkback = walkback + 1
	}

	for i := 0; i < walkback; i++ {
		if s[len(s)-i-1] != '0' && s[len(s)-i-1] != '.' {
			return prefix + s[:len(s)-i]
		}
	}

	return prefix + s[:len(s)-walkback]
}
