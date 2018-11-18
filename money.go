package money

import (
	"github.com/zenazn/money/decimal"
)

type Currency string

type Money struct {
	amt decimal.Decimal
}
