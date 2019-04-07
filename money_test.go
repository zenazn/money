package money

import (
	"testing"

	"github.com/zenazn/money/currency"
	"github.com/zenazn/money/decimal"
)

var constructorTests = []struct {
	v   Money
	amt string
	ccy currency.Currency
}{
	{Money{}, "0", nil},
	{Zero(currency.USD), "0", currency.USD},
	{New(decimal.FromI64(42), currency.EUR), "42", currency.EUR},
	{FromMinorCurrencyUnits(123, currency.USD), "1230000", currency.USD},
	{FromMinorCurrencyUnits(123, currency.JPY), "123000000", currency.JPY},
}

func TestConstructors(t *testing.T) {
	for i, test := range constructorTests {
		if s := test.v.amt.String(); s != test.amt {
			t.Errorf("[%d] amount expected %q got %q", i, test.amt, s)
		}
		if c := test.v.ccy; c != test.ccy {
			t.Errorf("[%d] amount expected %v got %v", i, test.ccy, c)
		}
	}
}

var stringTests = []struct {
	v Money
	s string
}{
	{Money{}, "0"},
	{Zero(currency.USD), "USD 0.00"},
	{FromMinorCurrencyUnits(123, currency.EUR), "EUR 1.23"},
	{FromMinorCurrencyUnits(23, currency.EUR), "EUR 0.23"},
	{FromMinorCurrencyUnits(3, currency.EUR), "EUR 0.03"},
	{FromMinorCurrencyUnits(12300, currency.EUR), "EUR 123.00"},
	{FromMinorCurrencyUnits(123, currency.EUR), "EUR 1.23"},
	{New(decimal.FromI64(42), currency.CAD), "CAD 0.000042"},
	{New(decimal.FromI64(12345678), currency.CAD), "CAD 12.345678"},
	{New(decimal.FromI64(12345600), currency.CAD), "CAD 12.3456"},
	{FromMinorCurrencyUnits(123, currency.JPY), "JPY 123"},
	{New(decimal.FromI64(12300000), currency.JPY), "JPY 12.3"},
	{New(decimal.FromI64(12000000), currency.JPY), "JPY 12"},
}

func TestString(t *testing.T) {
	for i, test := range stringTests {
		if s := test.v.String(); s != test.s {
			t.Errorf("[%d] expected %q got %q", i, test.s, s)
		}
	}
}
