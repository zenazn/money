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

// FromMinorUnits returns a Money with the given currency and amount, as
// represented as an integer number of minor currency units in that currency.
func FromMinorUnits(amt int64, ccy currency.Currency) Money {
	if ccy == nil {
		panic("money: no currency given")
	}
	sf := int64(1000000)
	u := ccy.Units()
	sfe := int(u.MajorUnitScalingFactorExponent - u.MinorUnitsInMajorUnitExponent)
	d := decimal.FromI64(amt)
	for i := 0; i < sfe; i++ {
		sf = sf * 10
		if sf > 1000000000000000 {
			d = d.Mul(decimal.NewRate(sf))
			sf = 1000000
		}
	}
	return Money{d.Mul(decimal.NewRate(sf)), ccy}
}

var ten = decimal.NewRate(10 * 1000000)

// Parse interprets the amount as a floating-point decimal string (using "." to
// separate the whole part from the fractional part, and without thousands
// separators or other adornments) and the currency as an ISO 4217 currency
// code, and returns a Money representing that value, or an error if either the
// amount or currency fails to parse.
func Parse(amt, ccy string) (Money, error) {
	c, err := currency.FromISOSymbol(ccy)
	if err != nil {
		return Money{}, err
	}

	sf := int(c.Units().MajorUnitScalingFactorExponent)
	d := decimal.Decimal{}
	dot := -1

	// Classic multiply-by-10-and-add parser. I'm sure we could make this
	// much faster if we wanted.
	for i, chr := range amt {
		if dot >= 0 && i-dot >= sf {
			return Money{}, fmt.Errorf("money: too precise")
		}

		if chr == '.' && i != 0 && dot == -1 {
			dot = i + 1
			continue
		} else if chr < '0' || chr > '9' {
			return Money{}, fmt.Errorf("money: bad char %q at position %d", chr, i)
		}
		d = d.Mul(ten).Add(decimal.FromI64(int64(chr - '0')))
	}

	if dot == len(amt) {
		// If we saw a dot at the very end, that's malformed
		return Money{}, fmt.Errorf("money: trailing dot")
	} else if dot == -1 {
		// If we never saw a dot, that's equivalent to it being at the end
		dot = len(amt)
	}

	// Scale out to the scale factor
	for i := len(amt) - dot; i < sf; i++ {
		d = d.Mul(ten)
	}

	return Money{d, c}, nil
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

// ComparableTo returns true if the two values are comparable (i.e., have the
// same currency).
func (m Money) ComparableTo(o Money) bool {
	return compat(m.ccy, o.ccy) == nil
}

func (m Money) compatCcy(o Money) currency.Currency {
	// Operations on currencyless zeroes are "sticky"
	if m.ccy != nil {
		return m.ccy
	}
	return o.ccy

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
