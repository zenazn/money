// Package currency defines a type for currency, as well as helpers for common
// ISO currencies.
package currency

import "errors"

// Units describes the relationship between one major unit of a currency and its
// various subdivisions.
type Units struct {
	// Ten raised to this power is the number of minor units in one major
	// unit of the currency.
	MinorUnitsInMajorUnitExponent uint8

	// Ten raised to this power is the number of minimum-representable-units
	// in one major currency unit. This is ideally a number large enough
	// that no reasonable computation in that currency would require more
	// precision than this, but which leaves enough digits to represent
	// every reasonable value. Unless you have a good reason to select
	// another value, you should select 6.
	//
	// It is common to store currency values as an integer number of
	// minimum-representable-units. Therefore, it is critical to never
	// change this value for an existing currency, lest previously-stored
	// values be interpreted with a different scaling factor.
	MajorUnitScalingFactorExponent uint8
}

// Currency is a type representing the unit in a system of money.
type Currency interface {
	// Symbol is an identifier that uniquely identifies a currency. For
	// government-issued currencies in common use, Symbol is defined to be
	// that currency's ISO 4217 currency code.
	Symbol() string
	// Units describes the relationship between one major unit of a currency
	// and its various subdivisions.
	Units() Units
}

type iso uint16

func (i iso) Symbol() string {
	if int(i)*3 > len(numberToSymbol) {
		panic("currency: iso currency out of range")
	}
	s := numberToSymbol[int(i)*3 : int(i)*3+3]
	if s == "..." {
		panic("currency: unknown iso currency")
	}
	return s
}

func (i iso) Units() Units {
	if int(i) > len(numberToSymbol) {
		panic("currency: iso currency out of range")
	}
	m := numberToMinor[int(i)]
	if m == 0xff {
		panic("currency: unknown iso currency")
	}
	return Units{m, 6}
}

var _ Currency = iso(0)

var NoSuchCurrency = errors.New("currency: no such ISO currency")

// FromISOSymbol returns a Currency object corresponding to the given ISO
// currency code, or NoSuchCurrency if no such currency exists.
func FromISOSymbol(s string) (Currency, error) {
	if c, ok := symbolToCurrency[s]; ok {
		return c, nil
	}
	return nil, NoSuchCurrency
}
