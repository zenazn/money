package money

import "github.com/zenazn/money/decimal"

// AddErr adds the two values and returns the result, or an error if the two
// values were incompatible.
func (m Money) AddErr(o Money) (Money, error) {
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

// SubErr subtracts the second value from the first and returns the result, or
// an error if the two values are incompatible.
func (m Money) SubErr(o Money) (Money, error) {
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

// Mul multiplies the value by the given scalar rate and returns the result.
func (m Money) Mul(r decimal.Rate) Money {
	return Money{m.amt.Mul(r), m.ccy}
}

// Div divides the value by the given scalar rate and returns the result.
func (m Money) Div(r decimal.Rate) Money {
	return Money{m.amt.Div(r), m.ccy}
}

// Neg negates the value and returns the result.
func (m Money) Neg() Money {
	return Money{m.amt.Neg(), m.ccy}
}

func (m Money) RoundToMinorUnits() Money {
	// TODO
	return m
}
